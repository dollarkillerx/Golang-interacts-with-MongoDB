/**
 * @Author: DollarKiller
 * @Description:
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 19:25 2019-10-27
 */
package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

type Book struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string
	Category string
	Weight   int
	Author   AuthorInfo
}

type AuthorInfo struct {
	Name    string
	Country string
}

const (
	categoryComputer = "计算机"
	categorySciFi    = "科幻"
	countryChina     = "中国"
	countryAmerica   = "美国"
)

var (
	books = []interface{}{
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "深入理解计算机操作系统",
			Category: categoryComputer,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "兰德尔 E.布莱斯特",
				Country: countryAmerica,
			},
		},
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "深入理解Linux内核",
			Category: categoryComputer,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "博韦，西斯特",
				Country: countryAmerica,
			},
		},
		&Book{
			Id:       primitive.NewObjectID(),
			Name:     "三体",
			Category: categorySciFi,
			Weight:   1,
			Author: AuthorInfo{
				Name:    "刘慈欣",
				Country: countryChina,
			},
		},
	}
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// 链接数据库
	client, e := mongo.Connect(context.Background(), opts)
	if e != nil {
		log.Fatal(e)
	}

	// 判断服务是否可用
	if e = client.Ping(context.TODO(), nil); e != nil {
		log.Fatal(e)
	}

	// 获取数据库和集合
	collection := client.Database("mydb").Collection("BOOK")

	// 清空数据库和集合
	e = collection.Drop(context.TODO())

	// 设置索引
	idx := mongo.IndexModel{
		Keys:    bsonx.Doc{{"name", bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	}

	idxRet, err := collection.Indexes().CreateOne(context.TODO(), idx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Collection.Indexes().CreateOne:", idxRet)

	// 插入一条数据
	insertOneResult, err := collection.InsertOne(context.TODO(), books[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Collection.InsertOne:", insertOneResult.InsertedID)

	// 插入多条数据
	insertManyResult, err := collection.InsertMany(context.TODO(), books[1:])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.InsertMany:", insertManyResult.InsertedIDs)

	// 获取数据总数
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(count)
	}
	log.Println("collection.CountDocuments:", count)

	// 查询单挑数据
	var one Book
	e = collection.FindOne(context.Background(), bson.M{"name": "三体"}).Decode(&one)
	if e != nil {
		log.Fatal(err)
	}
	log.Println("collection.FindOne: ", one)

	// 查询多条数据 方式一
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	var all []*Book
	err = cur.All(context.TODO(), &all)
	if err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())

	log.Fatalln("collection.Find curl.All", all)
	for _, one := range all {
		log.Println(one)
	}

	// 查询多条数据 方式二
	cur, err = collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find cur.Next:", b)
	}
	cur.Close(context.TODO())

	// 模糊查询
	cur, err = collection.Find(context.TODO(), bson.M{"name": primitive.Regex{Pattern: "深入"}})
	if err != nil {
		log.Fatal(err)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find name=primitive.Regex{深入}: ", b)
	}
	cur.Close(context.TODO())

	// 二级结构体查询
	// 二级结构体查询
	cur, err = collection.Find(context.Background(), bson.M{"author.country": countryChina})
	// cur, err = collection.Find(context.Background(), bson.D{bson.E{"author.country", countryChina}})
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.Background()) {
		var b Book
		if err = cur.Decode(&b); err != nil {
			log.Fatal(err)
		}
		log.Println("collection.Find author.country=", countryChina, ":", b)
	}
	cur.Close(context.Background())

	// 修改一条数据
	b1 := books[0].(*Book)
	b1.Weight = 2
	update := bson.M{"$set": b1}
	updateResult, err := collection.UpdateOne(context.Background(), bson.M{"name": b1.Name}, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.UpdateOne:", updateResult)

	// 修改一条数据，如果不存在则插入
	new := &Book{
		Id:       primitive.NewObjectID(),
		Name:     "球状闪电",
		Category: categorySciFi,
		Author: AuthorInfo{
			Name:    "刘慈欣",
			Country: countryChina,
		},
	}
	update = bson.M{"$set": new}
	updateOpts := options.Update().SetUpsert(true)
	updateResult, err = collection.UpdateOne(context.Background(), bson.M{"_id": new.Id}, update, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.UpdateOne:", updateResult)

	// 删除一条数据
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": new.Id})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.DeleteOne:", deleteResult)
}
