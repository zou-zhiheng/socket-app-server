package service

import (
	"app-server/global"
	"app-server/model"
	"app-server/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

func GetApi(name, currPage, pageSize, startTime, endTime string) utils.Response {
	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		return utils.ErrorMess("行数字段过长", err)
	}
	curr, err := strconv.ParseInt(currPage, 10, 64)
	if err != nil {
		return utils.ErrorMess("指定页面字段过长", err)
	}

	skip := (curr - 1) * size
	opt := options.FindOptions{
		Limit: &size,
		Skip:  &skip,
		Sort:  bson.M{"_id": -1},
	}

	filter := bson.M{}

	if startTime != "" && endTime != "" {
		filter["createTime"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	filter["name"] = bson.M{"$regex": name}

	count, _ := global.ApiColl.CountDocuments(context.TODO(), bson.M{})
	res, err := global.ApiColl.Find(context.TODO(), filter, &opt)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	var apiDB []model.Api
	if err = res.All(context.TODO(), &apiDB); err != nil {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", struct {
		Count int64       `json:"count" bson:"count"`
		Data  []model.Api `json:"data" bson:"data"`
	}{
		Count: count,
		Data:  apiDB,
	})
}

func CreateApi(api model.Api) utils.Response {

	if api.Name == "" || api.Path == "" || api.Method != "POST" && api.Method != "PUT" && api.Method != "DELETE" && api.Method != "GET"    {
		return utils.ErrorMess("失败", "名称,接口,方法不能为空")
	}

	filter := bson.M{
		"$or": []bson.M{
			bson.M{"name": api.Name},
			bson.M{"path": api.Path},
		},
	}

	if err := global.ApiColl.FindOne(context.TODO(), filter).Decode(&bson.M{}); err == nil {
		return utils.ErrorMess("失败,此api已存在", err)
	}

	res, err := global.ApiColl.InsertOne(context.TODO(), api)
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	return utils.ErrorMess("成功", res.InsertedID)
}

func UpdateApi(api model.Api) utils.Response {

	if api.Id.IsZero() {
		return utils.ErrorMess("失败,禁止访问!", nil)
	}

	if api.Name == "" || api.Path == "" || api.Method != "POST" && api.Method != "PUT" && api.Method != "DELETE" && api.Method != "GET"    {
		return utils.ErrorMess("失败", "名称,接口,方法不能为空")
	}

	var apiDB model.Device
	err := global.ApiColl.FindOne(context.TODO(), bson.M{"_id": api.Id}).Decode(&apiDB)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	res, err := global.ApiColl.UpdateOne(context.TODO(), bson.M{"_id": api.Id}, bson.M{"$set": api})
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	return utils.SuccessMess("成功", res)
}

func DeleteApi(idStr string) utils.Response {

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	//先删除对应用户的访问此api的权限
	_, err = global.UserColl.UpdateMany(context.TODO(), bson.M{"$where":"this.auth!=null"}, bson.M{"$pull": bson.M{"auth": id}})
	if err != nil {
		return utils.ErrorMess("用户更新", err)
	}

	//api删除
	res, err := global.ApiColl.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return utils.ErrorMess("失败", nil)
	}

	return utils.SuccessMess("成功", res.DeletedCount)
}
