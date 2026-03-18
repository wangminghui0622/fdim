package svc

import (
	"context"
	"time"

	"fdim/api/internal/config"
	"fdim/pkg/database"
	"fdim/pkg/grpcinterceptor"
	"fdim/pkg/model"
	"fdim/protocol/admin"
	"fdim/protocol/auth"
	"fdim/protocol/chat"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/msggateway"
	"fdim/protocol/third"
	"fdim/protocol/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext struct {
	Config             config.Config
	UserClient         user.UserClient
	FriendClient       user.FriendClient
	GroupClient        user.GroupClient
	AuthClient         auth.AuthClient
	ConversationClient conversation.ConversationClient
	MsgClient          msg.MsgClient
	ThirdClient        third.ThirdClient
	MsgGatewayClient   msggateway.MsgGatewayClient
	ChatClient         chat.ChatClient
	AdminClient        admin.AdminClient
	FavoriteDB         model.FavoriteMsgModel
	UserDB             *database.UserDatabase
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 RPC 客户端，添加拦截器以传递context中的用户信息
	userClient := user.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	friendClient := user.NewFriendClient(zrpc.MustNewClient(c.FriendRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	groupClient := user.NewGroupClient(zrpc.MustNewClient(c.GroupRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	authClient := auth.NewAuthClient(zrpc.MustNewClient(c.AuthRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	conversationClient := conversation.NewConversationClient(zrpc.MustNewClient(c.ConversationRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	msgClient := msg.NewMsgClient(zrpc.MustNewClient(c.MsgRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	thirdClient := third.NewThirdClient(zrpc.MustNewClient(c.ThirdRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
	
	var msgGatewayClient msggateway.MsgGatewayClient
	if len(c.MessageGatewayRpc.Etcd.Hosts) > 0 {
		msgGatewayClient = msggateway.NewMsgGatewayClient(zrpc.MustNewClient(c.MessageGatewayRpc).Conn())
	}

	var chatClient chat.ChatClient
	if c.ChatRpc.Etcd.Key != "" {
		chatClient = chat.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
		logx.Info("ChatRpc client connected")
	}

	var adminClient admin.AdminClient
	if c.AdminRpc.Etcd.Key != "" {
		adminClient = admin.NewAdminClient(zrpc.MustNewClient(c.AdminRpc, zrpc.WithDialOption(grpcinterceptor.GrpcClient())).Conn())
		logx.Info("AdminRpc client connected")
	}

	var favoriteDB model.FavoriteMsgModel
	var userDB *database.UserDatabase
	if c.Mongo.URI != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Mongo.URI))
		if err != nil {
			logx.Errorf("Failed to connect to MongoDB: %v", err)
		} else {
			db := mongoClient.Database(c.Mongo.Database)
			favoriteDB = database.NewFavoriteMsgMongo(db)
			mongoDB := database.NewMongoDBFromClient(mongoClient, c.Mongo.Database)
			userDB = database.NewUserDatabase(mongoDB)
			logx.Info("MongoDB connected for favorite and user search feature")
		}
	}

	return &ServiceContext{
		Config:             c,
		UserClient:         userClient,
		FriendClient:       friendClient,
		GroupClient:        groupClient,
		AuthClient:         authClient,
		ConversationClient: conversationClient,
		MsgClient:          msgClient,
		ThirdClient:        thirdClient,
		MsgGatewayClient:   msgGatewayClient,
		ChatClient:         chatClient,
		AdminClient:        adminClient,
		FavoriteDB:         favoriteDB,
		UserDB:             userDB,
	}
}
