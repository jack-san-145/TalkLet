package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

func Check_Private_MessagePartition(rdb *redis.Client, pool *pgxpool.Pool) {

	var Partition_exists bool
	today := time.Now().Day()
	fmt.Println("today day in month - ", today)

	//for checking only one time in the month
	if today != 1 { //change it to 1 to make the partition from day 1 of the month
		return
	}

	err := rdb.Get(context.Background(), "is_private_partition_created").Err()

	if err == redis.Nil {
		start_date := time.Now().Format("2006-01-02")
		after_one_month := time.Now().AddDate(0, 1, 0) // add exactly 1 month from today ex: 2005-12-30 -> 2006-01-30
		end_date := after_one_month.Format("2006-01-02")

		partition_table_name := "private_messages_" + start_date
		fmt.Println("private partition table name - ", partition_table_name)

		//checking if the partition table exist or not
		query_to_check_existence_of_partition := fmt.Sprintf(`select exists(select from pg_tables where schemaname='public' and tablename='%s')`, partition_table_name)
		pool.QueryRow(context.Background(), query_to_check_existence_of_partition).Scan(&Partition_exists)

		if Partition_exists {
			return
		}

		//query to create the partition table "private_messages_2025-08-12"
		query := fmt.Sprintf(`create table if not exists "%s" partition of all_private_messages for values from ('%s') to ('%s')`, partition_table_name, start_date, end_date)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating private partition - ", err)
			return
		}

		// query to set the primary key to that partition table "private_messages_2025-08-12"

		// alter table "private_ messages_2025-08-12" add constraint "pk_private_messages_2025-08-12" primary key(msg_id);
		primary_key := "pk_" + partition_table_name
		query = fmt.Sprintf(`alter table "%s" add constraint "%s" primary key(msg_id)`, partition_table_name, primary_key)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while altering the primary key - ", err)
			return
		}

		//create index "idx_private_sr_messages_2025-08-12" on "private_messages_2025-08-12" (sender_id,receiver_id);
		index_sender_receiver := "idx_private_sr_" + partition_table_name
		query = fmt.Sprintf(`create index if not exists "%s" on "%s" (sender_id,receiver_id);`, index_sender_receiver, partition_table_name)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating index for private partition table - ", err)
			return
		}

		//set the partition status in redis for 30 days
		err = rdb.Set(context.Background(), "is_private_partition_created", true, time.Hour*30*24).Err()
		if err != nil {
			fmt.Println("error while setting the private partition status to redis - ", err)
		}

	} else if err != nil {
		fmt.Println("error while accessing the private partition status - ", err)
	}

}

func Check_Group_MessagePartition(rdb *redis.Client, pool *pgxpool.Pool, dept string) {

	dept_group_msg_table := dept + "_group_all_messages"

	var Group_Partition_exists bool
	today := time.Now().Day()
	fmt.Println("today day in month - ", today)

	//for checking only one time in the month
	if today != 1 { //change it to 1 to make the partition from day 1 of the month
		return
	}

	redis_group_partition_status := "is_" + dept + "_group_partition_created" // is_cs_group_partition_created
	err := rdb.Get(context.Background(), redis_group_partition_status).Err()

	if err == redis.Nil {
		start_date := time.Now().Format("2006-01-02")
		after_one_month := time.Now().AddDate(0, 1, 0) // add exactly 1 month from today ex: 2005-12-30 -> 2006-01-30
		end_date := after_one_month.Format("2006-01-02")

		partition_table_name := dept + "_group_messages_" + start_date
		fmt.Println("group partition table name - ", partition_table_name)

		//checking if the group partition table exist or not
		query_to_check_existence_of_partition := fmt.Sprintf(`select exists(select from pg_tables where schemaname='public' and tablename='%s')`, partition_table_name)
		pool.QueryRow(context.Background(), query_to_check_existence_of_partition).Scan(&Group_Partition_exists)

		if Group_Partition_exists {
			return
		}

		//query to create the partition table "cs_group_messages_2025-08-12"
		query := fmt.Sprintf(`create table if not exists "%s" partition of %s for values from ('%s') to ('%s')`, partition_table_name, dept_group_msg_table, start_date, end_date)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating group partition - ", err)
			return
		}

		// query to set the primary key to that partition table "cs_group_messages_2025-08-12"

		// alter table "cs_group_messages_2025-08-12" add constraint "pk_cs_group_messages_2025-08-12" primary key(msg_id);
		primary_key := "pk_" + partition_table_name
		query = fmt.Sprintf(`alter table "%s" add constraint "%s" primary key(msg_id)`, partition_table_name, primary_key)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while altering the group primary key - ", err)
			return
		}

		//create index "idx_cs_group_sg_messages_2025-08-12" on "cs_group_messages_2025-08-12" (sender_id,receiver_id);
		index_sender_groupId := "idx_" + dept + "_group_sg_" + partition_table_name
		query = fmt.Sprintf(`create index if not exists "%s" on "%s" (sender_id,group_id);`, index_sender_groupId, partition_table_name)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating index for group partition table - ", err)
			return
		}

		//set the partition status in redis for 30 days
		err = rdb.Set(context.Background(), redis_group_partition_status, true, time.Hour*30*24).Err()
		if err != nil {
			fmt.Println("error while setting the group partition status to redis - ", err)
		}

	} else if err != nil {
		fmt.Println("error while accessing the group partition status - ", err)
	}

}
