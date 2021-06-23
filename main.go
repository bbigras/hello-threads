package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alecthomas/jsonschema"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/textileio/go-threads/api/client"
	"github.com/textileio/go-threads/core/thread"
	db2 "github.com/textileio/go-threads/db"
	"google.golang.org/grpc"
)

type Person struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	CreatedAt int    `json:"created_at"`
}

func getThread(name string, db0 *client.Client, dbs map[thread.ID]db2.Info) *thread.ID {
	for threadID := range dbs {

		dbInfo, errInfo := db0.GetDBInfo(context.Background(), threadID)
		if errInfo != nil {
			panic(errInfo)
		}

		if dbInfo.Name == name {
			return &threadID
		}
	}
	return nil
}

func main() {
	logging.SetupLogging(logging.Config{
		Format: logging.ColorizedOutput,
		Stdout: true,
		Level:  logging.LevelDebug,
	})

	var nFlag = flag.Int("c", 1, "computer number (1 or 2)")
	flag.Parse()

	// CLIENT
	db, err := client.NewClient("127.0.0.1:6006", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// Generate private key and save to file
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

	_, err3 := db.GetToken(context.Background(), myIdentity)
	if err3 != nil {
		panic(err3)
	}

	// computer 1
	if *nFlag == 1 {
		dbs, errDbs := db.ListDBs(context.Background())
		if errDbs != nil {
			panic(errDbs)
		}

		threadID := getThread("my-database", db, dbs)
		fmt.Println("threadID", threadID)

		if threadID == nil {
			fmt.Println("create database")
			threadID2 := thread.NewIDV1(thread.Raw, 32)
			options := db2.WithNewManagedName("my-database")
			err2 := db.NewDB(context.Background(), threadID2, options)
			if err2 != nil {
				panic(err2)
			}
			threadID = &threadID2
		}

		collection_exists := false
		collections, errCol := db.ListCollections(context.Background(), *threadID)
		if errCol != nil {
			panic(errCol)
		}
		for one := range collections {
			fmt.Println("collection> ", collections[one].Name)
			if collections[one].Name == "Persons" {
				collection_exists = true
			}
		}

		if !collection_exists {
			fmt.Println("create collection")
			reflector := jsonschema.Reflector{}
			mySchema := reflector.Reflect(&Person{}) // Generate a JSON Schema from a struct

			err15 := db.NewCollection(context.Background(), *threadID, db2.CollectionConfig{
				Name:   "Persons",
				Schema: mySchema,
				Indexes: []db2.Index{{
					Path:   "name", // Value matches json tags
					Unique: true,   // Create a unique index on "name"
				}},
			})
			if err15 != nil {
				panic(err15)
			}
		}

		query := db2.Where("name").Eq("Alice")
		results, err16 := db.Find(context.Background(), *threadID, "Persons", query, &Person{})
		if err16 != nil {
			panic(err16)
		}

		alice := results.([]*Person)

		if len(alice) == 0 {
			fmt.Println("add Alice")
			alice := &Person{
				ID:        "",
				Name:      "Alice",
				Age:       30,
				CreatedAt: int(time.Now().UnixNano()),
			}

			ids, err99 := db.Create(context.Background(), *threadID, "Persons", client.Instances{alice})
			if err99 != nil {
				panic(err99)
			}
			fmt.Println("> Success!", ids[0])
		} else {
			fmt.Println("alice", alice[0])
		}

		// info for computer 2
		fmt.Println("\n\nINFO FOR COMPUTER 2\n")
		dbs, errDbs2 := db.ListDBs(context.Background())
		if errDbs2 != nil {
			panic(errDbs2)
		}

		for threadID := range dbs {
			fmt.Println("thread", threadID)

			dbInfo, errInfo := db.GetDBInfo(context.Background(), threadID)
			if errInfo != nil {
				panic(errInfo)
			}
			fmt.Println("addrs", dbInfo.Addrs)
			fmt.Println("key", dbInfo.Key)
		}
	} else {
		// computer 2
		addr, err22 := ma.NewMultiaddr("/ip4/100.67.149.7/tcp/4006/p2p/12D3KooWEjMr3DccgP4tvUvpLEgfi7XEDDF1pxQ9dZimWLQPS67c/thread/bafk336kt326667ejcrpupk772wyfpjoik35wicwqyyn77d3fabnzdly")
		if err22 != nil {
			panic(err22)
		}

		key, err21 := thread.KeyFromString("bzfv6hms2bu4pqeedo333x3klnqcisszzqavq4mdquvm5sjqmwlvwtfhxe746yaoskqnxxa5hh3wwtm36otd3mkihwyteb7qe6ptgnky")
		if err21 != nil {
			panic(err21)
		}

		err20 := db.NewDBFromAddr(context.Background(), addr, key)
		if err20 != nil {
			panic(err20)
		}
	}
}
