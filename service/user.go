package service

import (
	"app-server/global"
	"app-server/middleware"
	"app-server/model"
	"app-server/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

func GetUserData(name, currPage, pageSize, startTime, endTime string) utils.Response {
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

	filter["isValid"] = true
	filter["name"] = bson.M{"$regex": name}

	count, _ := global.UserColl.CountDocuments(context.TODO(), bson.M{"isValid": true})
	res, err := global.UserColl.Find(context.TODO(), filter, &opt)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	var userDB []model.User
	if err = res.All(context.TODO(), &userDB); err != nil {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", struct {
		Count int64        `json:"count" bson:"count"`
		Data  []model.User `json:"data" bson:"data"`
	}{
		Count: count,
		Data:  userDB,
	})
}

func CreateUser(user model.User) utils.Response {

	if user.Name == "" || user.Account == "" || user.Password == "" {
		return utils.ErrorMess("失败,参数错误", nil)
	}

	if err := global.UserColl.FindOne(context.TODO(), bson.M{"account": user.Account, "isValid": true}).Decode(&bson.M{}); err == nil {
		return utils.ErrorMess("失败,该用户已存在", err)
	}

	//密码加密
	//生成种子
	rand.Seed(time.Now().Unix())
	//生成盐
	user.Salt = strconv.FormatInt(rand.Int63(), 10)
	//密码加盐加密
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+user.Salt), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorMess("密码加密失败", err)
	}

	user.Password = string(encryptedPass)
	user.IsValid = true
	user.CreateTime = utils.TimeFormat(time.Now())
	//记录插入

	res, err := global.UserColl.InsertOne(context.TODO(), user)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", res.InsertedID)
}

func UpdateUser(user model.User, flag string) utils.Response {

	var userDB model.User
	if err := global.UserColl.FindOne(context.TODO(), bson.M{"_id": user.Id}).Decode(&userDB); err != nil {
		return utils.ErrorMess("该用户不存在", err)
	}

	user.UpdateTime = utils.TimeFormat(time.Now())
	//避免关键数据被修改
	user.Account = userDB.Account
	user.CreateTime = userDB.CreateTime
	user.IsValid = userDB.IsValid

	if flag == "true" {

		if user.Password == "" {
			return utils.ErrorMess("失败,密码不能为空", nil)
		}

		//密码加密
		//生成种子
		rand.Seed(time.Now().Unix())
		//生成盐
		user.Salt = strconv.FormatInt(rand.Int63(), 10)
		//密码加盐加密
		encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+user.Salt), bcrypt.DefaultCost)
		if err != nil {
			return utils.ErrorMess("密码加密失败", err)
		}

		user.Password = string(encryptedPass)
	} else {
		user.Password = userDB.Password
		user.Salt = userDB.Salt
	}

	res, err := global.UserColl.UpdateOne(context.TODO(), bson.M{"_id": user.Id}, bson.D{{"$set", user}})
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", res)
}

func DeleteUser(idStr string) utils.Response {

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return utils.ErrorMess("参数错误", nil)
	}

	err = global.UserColl.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&bson.M{})
	if err != nil {
		return utils.ErrorMess("该用户不存在", err)
	}

	res, err := global.UserColl.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"isValid": false}})
	if err != nil || res.ModifiedCount == 0 {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", nil)
}

func UserLogin(user model.User) utils.Response {

	if user.Account == "" {
		return utils.ErrorMess("参数错误", nil)
	}

	var userDB model.User
	err := global.UserColl.FindOne(context.TODO(), bson.M{"account": user.Account}).Decode(&userDB)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	//密码验证
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password+userDB.Salt))
	if err != nil {
		return utils.ErrorMess("密码错误", err)
	}

	token, err := middleware.CreateToken(userDB)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	data := map[string]interface{}{
		"_id":     userDB.Id,
		"name":    userDB.Name,
		"account": userDB.Account,
		"auth":    userDB.Auth,
		"token":   token,
	}

	return utils.SuccessMess("成功", data)
}
