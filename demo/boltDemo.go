package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte("bucketName1"))
		if b1 == nil {
			// b1通不存在，则创建
			b1, err = tx.CreateBucket([]byte("bucketName1"))
			if err != nil {
				log.Panic(err)
			}
		}
		// bucket创建完成，开始写入数据
		err = b1.Put([]byte("name1"), []byte("Lily"))
		if err != nil {
			fmt.Printf("写入数据失败name1: Lily\n")
		}

		err = b1.Put([]byte("name2"), []byte("Jim"))
		if err != nil {
			fmt.Printf("写入数据失败name2: Jim\n")
		}

		// 读取数据
		name1 := b1.Get([]byte("name1"))
		name2 := b1.Get([]byte("name2"))
		name3 := b1.Get([]byte("name3"))
		fmt.Printf("name1: %s\n", name1)
		fmt.Printf("name2: %s\n", name2)
		fmt.Printf("name3: %s\n", name3)

		return nil
	})
}
