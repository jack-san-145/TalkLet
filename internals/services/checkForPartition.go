package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

func CheckForMessagePartition(rdb *redis.Client, pool *pgxpool.Pool) {
	today := time.Now().Day()
	fmt.Println("today day in month - ", today)

	//for checking only one time in the month
	if today != 12 { //change it to 1 to make the partition from day 1 of the month
		return
	}

	err := rdb.Get(context.Background(), "is_partition_created").Err()

	if err == redis.Nil {
		start_date := time.Now().Format("2006-01-02")
		after_one_month := time.Now().AddDate(0, 1, 0) // add exactly 1 month from today ex: 2005-12-30 -> 2006-01-30
		end_date := after_one_month.Format("2006-01-02")

		partition_table_name := "messages_" + start_date
		fmt.Println("partition table name - ", partition_table_name)

		//query to create the partition table "messages_2025-08-12"
		query := fmt.Sprintf(`create table if not exists "%s" partition of all_messages for values from ('%s') to ('%s')`, partition_table_name, start_date, end_date)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating partition - ", err)
			return
		}

		// query to set the primary key to that partition table "messages_2025-08-12"

		// alter table "messages_2025-08-12" add constraint "pk_messages_2025-08-12" primary key(msg_id);
		primary_key := "pk_" + partition_table_name
		query = fmt.Sprintf(`alter table "%s" add constraint "%s" primary key(msg_id)`, partition_table_name, primary_key)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while altering the primary key - ", err)
			return
		}

		//create index "idx_sr_messages_2025-08-12" on "messages_2025-08-12" (sender_id,receiver_id);
		index_sender_receiver := "idx_sr_" + partition_table_name
		query = fmt.Sprintf(`create index if not exists "%s" on "%s" (sender_id,receiver_id);`, index_sender_receiver, partition_table_name)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while creating index for partition table - ", err)
			return
		}

		//set the partition status in redis
		err = rdb.Set(context.Background(), "is_partition_created", true, time.Hour*30*24).Err()
		if err != nil {
			fmt.Println("error while setting the partition status to redis - ", err)
		}

	} else if err != nil {
		fmt.Println("error while accessing the partition status - ", err)
	}

}
