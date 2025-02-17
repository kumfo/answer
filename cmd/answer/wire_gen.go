// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/segmentfault/answer/internal/base/conf"
	"github.com/segmentfault/answer/internal/base/data"
	"github.com/segmentfault/answer/internal/base/middleware"
	"github.com/segmentfault/answer/internal/base/server"
	"github.com/segmentfault/answer/internal/base/translator"
	"github.com/segmentfault/answer/internal/controller"
	"github.com/segmentfault/answer/internal/controller_backyard"
	"github.com/segmentfault/answer/internal/repo"
	"github.com/segmentfault/answer/internal/repo/activity"
	"github.com/segmentfault/answer/internal/repo/activity_common"
	"github.com/segmentfault/answer/internal/repo/auth"
	"github.com/segmentfault/answer/internal/repo/captcha"
	"github.com/segmentfault/answer/internal/repo/collection"
	"github.com/segmentfault/answer/internal/repo/comment"
	"github.com/segmentfault/answer/internal/repo/common"
	"github.com/segmentfault/answer/internal/repo/config"
	"github.com/segmentfault/answer/internal/repo/export"
	"github.com/segmentfault/answer/internal/repo/meta"
	"github.com/segmentfault/answer/internal/repo/notification"
	"github.com/segmentfault/answer/internal/repo/rank"
	"github.com/segmentfault/answer/internal/repo/reason"
	"github.com/segmentfault/answer/internal/repo/report"
	"github.com/segmentfault/answer/internal/repo/revision"
	"github.com/segmentfault/answer/internal/repo/tag"
	"github.com/segmentfault/answer/internal/repo/unique"
	"github.com/segmentfault/answer/internal/repo/user"
	"github.com/segmentfault/answer/internal/router"
	"github.com/segmentfault/answer/internal/service"
	"github.com/segmentfault/answer/internal/service/action"
	activity2 "github.com/segmentfault/answer/internal/service/activity"
	"github.com/segmentfault/answer/internal/service/answer_common"
	auth2 "github.com/segmentfault/answer/internal/service/auth"
	"github.com/segmentfault/answer/internal/service/collection_common"
	comment2 "github.com/segmentfault/answer/internal/service/comment"
	export2 "github.com/segmentfault/answer/internal/service/export"
	"github.com/segmentfault/answer/internal/service/follow"
	meta2 "github.com/segmentfault/answer/internal/service/meta"
	notification2 "github.com/segmentfault/answer/internal/service/notification"
	"github.com/segmentfault/answer/internal/service/notification_common"
	"github.com/segmentfault/answer/internal/service/object_info"
	"github.com/segmentfault/answer/internal/service/question_common"
	rank2 "github.com/segmentfault/answer/internal/service/rank"
	reason2 "github.com/segmentfault/answer/internal/service/reason"
	report2 "github.com/segmentfault/answer/internal/service/report"
	"github.com/segmentfault/answer/internal/service/report_backyard"
	"github.com/segmentfault/answer/internal/service/report_handle_backyard"
	"github.com/segmentfault/answer/internal/service/revision_common"
	"github.com/segmentfault/answer/internal/service/service_config"
	tag2 "github.com/segmentfault/answer/internal/service/tag"
	"github.com/segmentfault/answer/internal/service/tag_common"
	"github.com/segmentfault/answer/internal/service/uploader"
	"github.com/segmentfault/answer/internal/service/user_backyard"
	"github.com/segmentfault/answer/internal/service/user_common"
	"github.com/segmentfault/pacman"
	"github.com/segmentfault/pacman/log"
)

// Injectors from wire.go:

// initApplication init application.
func initApplication(debug bool, serverConf *conf.Server, dbConf *data.Database, cacheConf *data.CacheConf, i18nConf *translator.I18n, swaggerConf *router.SwaggerConfig, serviceConf *service_config.ServiceConfig, logConf log.Logger) (*pacman.Application, func(), error) {
	staticRouter := router.NewStaticRouter(serviceConf)
	i18nTranslator, err := translator.NewTranslator(i18nConf)
	if err != nil {
		return nil, nil, err
	}
	langController := controller.NewLangController(i18nTranslator)
	engine, err := data.NewDB(debug, dbConf)
	if err != nil {
		return nil, nil, err
	}
	cache, cleanup, err := data.NewCache(cacheConf)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup2, err := data.NewData(engine, cache)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authRepo := auth.NewAuthRepo(dataData)
	authService := auth2.NewAuthService(authRepo)
	configRepo := config.NewConfigRepo(dataData)
	userRepo := user.NewUserRepo(dataData, configRepo)
	uniqueIDRepo := unique.NewUniqueIDRepo(dataData)
	activityRepo := repo.NewActivityRepo(dataData, uniqueIDRepo, configRepo)
	userRankRepo := rank.NewUserRankRepo(dataData, configRepo)
	userActiveActivityRepo := activity.NewUserActiveActivityRepo(dataData, activityRepo, userRankRepo, configRepo)
	emailRepo := export.NewEmailRepo(dataData)
	emailService := export2.NewEmailService(configRepo, emailRepo)
	userService := service.NewUserService(userRepo, userActiveActivityRepo, emailService, authService, serviceConf)
	captchaRepo := captcha.NewCaptchaRepo(dataData)
	captchaService := action.NewCaptchaService(captchaRepo)
	uploaderService := uploader.NewUploaderService(serviceConf)
	userController := controller.NewUserController(authService, userService, captchaService, emailService, uploaderService)
	commentRepo := comment.NewCommentRepo(dataData, uniqueIDRepo)
	commentCommonRepo := comment.NewCommentCommonRepo(dataData, uniqueIDRepo)
	userCommon := usercommon.NewUserCommon(userRepo)
	answerRepo := repo.NewAnswerRepo(dataData, uniqueIDRepo, userRankRepo, activityRepo)
	questionRepo := repo.NewQuestionRepo(dataData, uniqueIDRepo)
	tagRepo := tag.NewTagRepo(dataData, uniqueIDRepo)
	objService := object_info.NewObjService(answerRepo, questionRepo, commentCommonRepo, tagRepo)
	voteRepo := activity_common.NewVoteRepo(dataData, activityRepo)
	commentService := comment2.NewCommentService(commentRepo, commentCommonRepo, userCommon, objService, voteRepo)
	rankService := rank2.NewRankService(userCommon, userRankRepo, objService, configRepo)
	commentController := controller.NewCommentController(commentService, rankService)
	reportRepo := report.NewReportRepo(dataData, uniqueIDRepo)
	reportService := report2.NewReportService(reportRepo, objService)
	reportController := controller.NewReportController(reportService, rankService)
	serviceVoteRepo := activity.NewVoteRepo(dataData, uniqueIDRepo, configRepo, activityRepo, userRankRepo, voteRepo)
	voteService := service.NewVoteService(serviceVoteRepo, uniqueIDRepo, configRepo, questionRepo, answerRepo, commentCommonRepo, objService)
	voteController := controller.NewVoteController(voteService)
	revisionRepo := revision.NewRevisionRepo(dataData, uniqueIDRepo)
	revisionService := revision_common.NewRevisionService(revisionRepo, userRepo)
	followRepo := activity_common.NewFollowRepo(dataData, uniqueIDRepo, activityRepo)
	tagService := tag2.NewTagService(tagRepo, revisionService, followRepo)
	tagController := controller.NewTagController(tagService, rankService)
	followFollowRepo := activity.NewFollowRepo(dataData, uniqueIDRepo, activityRepo)
	followService := follow.NewFollowService(followFollowRepo, followRepo, tagRepo)
	followController := controller.NewFollowController(followService)
	collectionRepo := collection.NewCollectionRepo(dataData, uniqueIDRepo)
	collectionGroupRepo := collection.NewCollectionGroupRepo(dataData)
	tagRelRepo := tag.NewTagListRepo(dataData)
	tagCommonService := tagcommon.NewTagCommonService(tagRepo, tagRelRepo, revisionService)
	collectionCommon := collectioncommon.NewCollectionCommon(collectionRepo)
	answerCommon := answercommon.NewAnswerCommon(answerRepo)
	metaRepo := meta.NewMetaRepo(dataData)
	metaService := meta2.NewMetaService(metaRepo)
	questionCommon := questioncommon.NewQuestionCommon(questionRepo, answerRepo, voteRepo, followRepo, tagCommonService, userCommon, collectionCommon, answerCommon, metaService, configRepo)
	collectionService := service.NewCollectionService(collectionRepo, collectionGroupRepo, questionCommon)
	collectionController := controller.NewCollectionController(collectionService)
	answerActivityRepo := activity.NewAnswerActivityRepo(dataData, activityRepo, userRankRepo)
	questionActivityRepo := activity.NewQuestionActivityRepo(dataData, activityRepo, userRankRepo)
	answerActivityService := activity2.NewAnswerActivityService(answerActivityRepo, questionActivityRepo)
	questionService := service.NewQuestionService(questionRepo, tagCommonService, questionCommon, userCommon, revisionService, metaService, collectionCommon, answerActivityService)
	questionController := controller.NewQuestionController(questionService, rankService)
	answerService := service.NewAnswerService(answerRepo, questionRepo, questionCommon, userCommon, collectionCommon, userRepo, revisionService, answerActivityService, answerCommon, voteRepo)
	answerController := controller.NewAnswerController(answerService, rankService)
	searchRepo := repo.NewSearchRepo(dataData, uniqueIDRepo, userCommon)
	searchService := service.NewSearchService(searchRepo, tagRepo, userCommon, followRepo)
	searchController := controller.NewSearchController(searchService)
	serviceRevisionService := service.NewRevisionService(revisionRepo, userCommon, questionCommon, answerService)
	revisionController := controller.NewRevisionController(serviceRevisionService)
	rankController := controller.NewRankController(rankService)
	commonRepo := common.NewCommonRepo(dataData, uniqueIDRepo)
	reportHandle := report_handle_backyard.NewReportHandle(questionCommon, commentRepo, configRepo)
	reportBackyardService := report_backyard.NewReportBackyardService(reportRepo, userCommon, commonRepo, answerRepo, questionRepo, commentCommonRepo, reportHandle, configRepo)
	controller_backyardReportController := controller_backyard.NewReportController(reportBackyardService)
	userBackyardRepo := user.NewUserBackyardRepo(dataData)
	userBackyardService := user_backyard.NewUserBackyardService(userBackyardRepo)
	userBackyardController := controller_backyard.NewUserBackyardController(userBackyardService)
	reasonRepo := reason.NewReasonRepo(configRepo)
	reasonService := reason2.NewReasonService(reasonRepo)
	reasonController := controller.NewReasonController(reasonService)
	themeController := controller_backyard.NewThemeController()
	siteInfoRepo := repo.NewSiteInfo(dataData)
	siteInfoService := service.NewSiteInfoService(siteInfoRepo)
	siteInfoController := controller_backyard.NewSiteInfoController(siteInfoService)
	siteinfoController := controller.NewSiteinfoController(siteInfoService)
	notificationRepo := notification.NewNotificationRepo(dataData)
	notificationCommon := notificationcommon.NewNotificationCommon(dataData, notificationRepo, userCommon, activityRepo, followRepo, objService)
	notificationService := notification2.NewNotificationService(dataData, notificationRepo, notificationCommon)
	notificationController := controller.NewNotificationController(notificationService)
	answerAPIRouter := router.NewAnswerAPIRouter(langController, userController, commentController, reportController, voteController, tagController, followController, collectionController, questionController, answerController, searchController, revisionController, rankController, controller_backyardReportController, userBackyardController, reasonController, themeController, siteInfoController, siteinfoController, notificationController)
	swaggerRouter := router.NewSwaggerRouter(swaggerConf)
	uiRouter := router.NewUIRouter()
	authUserMiddleware := middleware.NewAuthUserMiddleware(authService)
	ginEngine := server.NewHTTPServer(debug, staticRouter, answerAPIRouter, swaggerRouter, uiRouter, authUserMiddleware)
	application := newApplication(serverConf, ginEngine)
	return application, func() {
		cleanup2()
		cleanup()
	}, nil
}
