package demo2

// db.getCollection('storage').find({"create_time": {$gt: 1611872395}}).sort({"create_time":-1})
func TestM() {
    var fixoptions options.FindOptions
    fixoptions.Sort = bson.D{{"create_time", -1}}

    filter := bson.D{{"create_time", bson.M{
      "$gt": 1611872395,
    }}}

    find, err := collection.Find(context.TODO(), filter, &fixoptions)
    if err != nil {
      log.Fatalln(err)
    }
    defer find.Close(context.TODO())

    var tasks []models.Storage
    if err := find.All(context.TODO(), &tasks); err != nil {
      log.Fatalln(err)
    }
	
}
