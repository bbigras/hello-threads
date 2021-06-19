package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"flag"
	// "time"

	ma "github.com/multiformats/go-multiaddr"
	// crypto "github.com/libp2p/go-libp2p-crypto"
	// "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-threads/api/client"
	"github.com/textileio/go-threads/core/thread"

	// db2 "github.com/textileio/go-threads/db"
	// "github.com/alecthomas/jsonschema"
	"google.golang.org/grpc"
)

type Person struct {
    ID        string `json:"_id"`
    Name      string `json:"name"`
    Age       int    `json:"age"`
    CreatedAt int    `json:"created_at"`
}

func main() {
	db, err := client.NewClient("127.0.0.1:6006", grpc.WithInsecure())
	if err != nil {
        panic(err)
	}

	var nFlag = flag.Bool("info", false, "show current DB info")
	flag.Parse()

	if *nFlag {
		dbs, errDbs := db.ListDBs(context.Background())
		if errDbs != nil {
			panic(errDbs)
		}

		for threadID := range dbs {
			fmt.Println("thread", threadID)

			dbInfo, errInfo := db.GetDBInfo(context.Background(), threadID)
			if errInfo != nil {
				panic(errInfo)
			}
			fmt.Println("addrs", dbInfo.Addrs)
			fmt.Println("key", dbInfo.Key)

			collections, errCol := db.ListCollections(context.Background(), threadID)
			if errCol != nil {
				panic(errCol)
			}
			for one := range collections {
				fmt.Println("collection> ", one)
			}

		}
		panic("done")
	}

	_, err7 := os.Stat("privkey")
	if err7 != nil {
		fmt.Println("generating key")

		privateKey, _, err4 := crypto.GenerateEd25519Key(rand.Reader) // Private key is kept locally
		if err4 != nil {
			panic(err4)
		}

		key2, err6 := crypto.MarshalPrivateKey(privateKey)
		if err6 != nil {
			panic(err6)
		}

		err5 := ioutil.WriteFile("privkey", key2, 0600)
		if err5 != nil {
			panic(err5)
		}
	}

	data, err8 := ioutil.ReadFile("privkey")
	if err8 != nil {
		panic(err8)
	}

	privateKey, err9 := crypto.UnmarshalPrivateKey(data)
	if err9 != nil {
		panic(err9)
	}

	myIdentity := thread.NewLibp2pIdentity(privateKey)
	fmt.Println("identity", myIdentity.GetPublic())

	// threadToken, err3 := db.GetToken(context.Background(), myIdentity)
	_, err3 := db.GetToken(context.Background(), myIdentity)
	if err3 != nil {
        panic(err3)
		}
	/*

	thread_pub_key, err10 := threadToken.PubKey()
	if err10 != nil {
		panic(err10)
	}
	fmt.Println("thread token", thread_pub_key)

	_, err12 := os.Stat("threadid")
	if err12 != nil {
		fmt.Println("new thread")
		threadID := thread.NewIDV1(thread.Raw, 32)

		err11 := ioutil.WriteFile("threadid", threadID.Bytes(), 0600)
		if err11 != nil {
			panic(err11)
		}

		err2 := db.NewDB(context.Background(), threadID)
		if err2 != nil {
			panic(err2)
		}
	}

	data2, err13 := ioutil.ReadFile("threadid")
	if err13 != nil {
		panic(err13)
	}

	// threadID := thread.ID(string(data2))
	threadID := thread.ID(data2)
	*/


	// on another computer
	// dbInfo, err := db.GetDBInfo(context.Background(), threadID)
	// fmt.Println("addrs", dbInfo.Addrs)
	// fmt.Println("key", dbInfo.Key)
	// err14 := db.NewDBFromAddr(context.Background(), dbInfo.Addrs[0], dbInfo.Key)
	// if err14 != nil {
		// panic(err14)
	// }

	// fmt.Println("threadID", threadID)

	// addr, err22 := ma.NewMultiaddr("/ip4/172.27.0.3/tcp/4006/p2p/12D3KooWQZUHsPQzcyzydWLcabNwKb9zhvgsYSfBwDTB3cZi33kT/thread/bafktrg35keacmar4neyegxcinck77z3scjctjmipkpkovyzi6g52coq")

	// addr, err22 := ma.NewMultiaddr("/ip4/100.67.149.7/tcp/4006/p2p/12D3KooWQZUHsPQzcyzydWLcabNwKb9zhvgsYSfBwDTB3cZi33kT/thread/bafktrg35keacmar4neyegxcinck77z3scjctjmipkpkovyzi6g52coq")
	addr, err22 := ma.NewMultiaddr("/ip4/100.67.149.7/tcp/4006/p2p/12D3KooWQZUHsPQzcyzydWLcabNwKb9zhvgsYSfBwDTB3cZi33kT/thread/bafktrg35keacmar4neyegxcinck77z3scjctjmipkpkovyzi6g52coq")
	if err22 != nil {
		panic(err22)
	}

	key, err21 := thread.KeyFromString("bm66r4d64fbrsrzsjpyjn7pwqznr3zvxbbm3w5mtsr6mcn5lymi3oxtz5kklur4r274xvbbimqbcjzcdegnigeyujq5tdreonqkyhv3y")
	if err21 != nil {
		panic(err21)
	}

	err20 := db.NewDBFromAddr(context.Background(), addr, key)
	if err20 != nil {
		panic(err20)
	}


	// reflector := jsonschema.Reflector{}
	// mySchema := reflector.Reflect(&Person{}) // Generate a JSON Schema from a struct

	// err15 := db.NewCollection(context.Background(), threadID, db2.CollectionConfig{
	// 	Name:    "Persons",
	// 	Schema:  mySchema,
	// 	Indexes: []db2.Index{{
	// 		Path:   "name", // Value matches json tags
	// 		Unique: true, // Create a unique index on "name"
	// 	}},
	// })
	// if err15 != nil {
	// 	panic(err15)
	// }

	   /*
err := db.UpdateCollection(context.Background(), myThreadID, db.CollectionConfig{
    Name:    "Persons",
    Schema:  mySchema,
    Indexes: []db.Index{{
        Path:   "name",
        Unique: true,
    },
    {
        Path: "created_at", // Add an additional index on "created_at"
    }},
	})
	*/

	/*
	alice := &Person{
		ID:        "",
		Name:      "Alice",
		Age:       30,
		CreatedAt: int(time.Now().UnixNano()),
	}

	ids, err := db.Create(context.Background(), threadID, "Persons", client.Instances{alice})
	fmt.Println("> Success!", ids[0])
	*/


	// query := db2.Where("name").Eq("Alice")
	// results, err16 := db.Find(context.Background(), threadID, "Persons", query, &Person{})
	// if err16 != nil {
	// 	panic(err16)
	// }


	// alice := results.([]*Person)

	// fmt.Println("len", len(alice))

	// if len(alice) > 0 {
	// 	fmt.Println("> Success!", alice[0])
	// }
}
