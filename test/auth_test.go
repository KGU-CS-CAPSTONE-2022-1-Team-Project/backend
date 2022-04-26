package test

import (
	"backend/internal/auth"
	"backend/internal/auth/dao"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func init() {
	viper.SetConfigName("client_secret")
	viper.SetConfigType("json")
	viper.SetConfigType("")
	viper.AddConfigPath("../configs/auth")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper error: %v", err))
	}
}
func dbconnection() *gorm.DB {
	config := viper.GetStringMapString("db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"], config["password"], config["host"], config["port"], config["main_db"])
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestCreateUser(t *testing.T) {
	user := auth.User{
		Email: "test@test.com",
	}
	userDB := dao.User{
		Email: user.Email,
	}
	defer func() {
		dbconnection().Unscoped().
			Delete(&userDB)
	}()
	err := userDB.Create()
	require.Nil(t, err, "생성 실패")
	result, err := userDB.ReadByID()
	require.Nil(t, err, "조회 실패")
	assert.Equal(t, user.Email, result.Email, "잘못된 입력")
}

func TestMigration(t *testing.T) {
	user := dao.User{}
	err := user.Migration()
	require.Nil(t, err, "마이그레이션")
}
