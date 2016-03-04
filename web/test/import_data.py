#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json

HOST = "http://127.0.0.1:8080"


APPS = ["1P_1080Wallpapers", "1P_47inchWallpapers", "1P_BigWallpapers", "1P_FineWallpapers", "1P_GlassWallpapers", "1P_GrassWallpapers", "1P_HugeWallpapers", "1P_LargeWallpapers", "1P_LockWallpapers", "1P_PanoramaWallpaperWallpapers", "1P_PhabletWallpapers", "1P_RainbowWallpapers", "1P_SapphireWallpapers", "1P_SmartWallpapers", "1P_ThemeWallpapers", "1P_Wallpapers", "1P_Wallpapers8", "1P_WallpapersHD", "1P_iOS8Wallpapers", "1P_iPhone6Wallpapers", "4P_CoolColorWallpapers", "4P_FunnyPix", "4P_FunnyPix_iPad", "4P_HarryPotterWallpapers", "4P_HelloKittyWallpapers", "4P_JustinBieberWallpapers", "4P_LockScreenWallpapers", "4P_LoveQuotesWallpapers", "4P_QuotesCNWallpapers", "4P_QuotesWallpapers", "4P_QuotesWallpapersPro", "4P_WallpaperBoxHD", "6P_ChristmasMagazineWallpapers", "6P_CleanWallpapers", "6P_ColorMagazineWallpapers", "6P_FlickrPics", "6P_HelloKittyMagazineWallpapers", "6P_HotPics", "6P_InstagramPics", "6P_LoveMagazineWallpapers", "6P_LoveMagazineWallpapers_iPad", "6P_MagazineWallpapers_iPad", "6P_MovieMagazineWallpapers", "6P_NewWallpapers", "6P_NewWallpapers_iPad", "6P_OneDirectionMagazineWallpapers", "6P_SnowMagazineWallpapers", "9P_AnimalWallpapers", "9P_ChristmasBGWallpapers", "9P_ChristmasCHWallpapers", "9P_ChristmasWallpapers", "9P_FashionWallpapers", "9P_HDRWallpapers", "9P_Halloween2012Wallpapers", "9P_HalloweenWallpapers", "9P_IconSkinsWallpapers", "9P_KidWallpapers", "9P_KidWallpapers_iPad", "9P_LolsotrueWallpapers", "9P_LoveQuotesWallpapers", "9P_MilitaryWallpapers", "9P_NewYearWallpapers", "9P_NewYearWallpapers_iPad", "9P_QuotesWallpapers", "9P_RetinaWallpapers", "9P_RetinaWallpapers_iPad", "9P_SeasonsWallpapers", "9P_Shelf1136Wallpapers", "9P_ShelfWallpapers", "9P_ShelfWallpapersPro", "9P_WallpapersBox", "9P_iOS7OnlyWallpaper", "9P_iOS7OnlyWallpaper_iPad", "9P_iOS7Wallpapershd", "9P_iOS7Wallpapershd_iPad", "9P_iPhone5Wallpapers", "9P_iPhone5Wallpapers_iPad", "AR_AppBox", "AR_AppGreen", "AR_AppRank", "AR_AppRed", "AR_COC", "AR_Common", "AR_CoverFree", "AR_CoverFreeCN", "AR_Desktop", "AR_Explore", "AR_FreeAppToday", "AR_GameFreeCN", "AR_GameLimitedFree", "AR_GameVideo", "AR_HDFreeCN", "AR_IconFreeCN", "AR_IconFreeCN_3_7", "AR_IconFreeCN_3_8", "AR_IconSortLimitedFree", "AR_LOL", "AR_Launcher", "AR_LimitedFree", "AR_ListSortLimitedFree", "AR_MY", "AR_MustAppsCN", "AR_SpeedU1", "AR_WePlay", "AR_XMStatusCN", "AR_XYPlus", "AR_XianYou", "AR_XianYouAnswer", "AR_XianYouAnswerHD", "AR_XianYouApp", "AR_XianYouAsk", "AR_XianYouGame", "AR_XianYouGameNews", "AR_XianYouVideo", "AR_iOSBooklet", "AR_iOSHandbook", "AR_iOSManual", "AT_IconER", "AT_IconERLite", "AT_IconERPro", "AV_AutoVideo", "AV_BeautyGirl", "AV_FunnyTime", "AY_Ipomelo", "AY_PhoneCase", "AY_Shoujike", "ChristmasRingtones", "Downloader", "FBPhoto", "FunnyRingtones", "GlowDIYBackgrounds", "Husky", "InstaMagic", "JokeAppEn", "POP_Ringtones", "PicsForTest", "Podcasts_News", "PopEmotion", "RT_360Ringtones", "RT_365Ringtones", "RT_EverRingtones", "RT_KuaiBo", "RT_PerfectRingtones", "RT_PowerRingtones", "RT_RingtonesBox", "RT_RingtonesBoxPro", "RT_RingtonesBox_Android", "RT_RingtonesStyle", "RT_iOS8Ringtones", "RT_iPhone6PlusRingtones", "RT_iPhone6Ringtones", "Reddit_AdultSecrets", "Reddit_AdviceCn", "Reddit_AdviceEn", "Reddit_AppleSchool", "Reddit_BabyCnA", "Reddit_ChatBoardEn", "Reddit_CoolFactsEn", "Reddit_DailyJokesenA", "Reddit_DirtyJokes", "Reddit_EmojiForum", "Reddit_EnglishquotescnA", "Reddit_EntertainmentNewsCn", "Reddit_EpicFail", "Reddit_FestivalCN", "Reddit_ForTestOnly", "Reddit_ForTestOnlyA", "Reddit_FunnyPicsCn", "Reddit_FunnyPicsEn", "Reddit_GamePark", "Reddit_GameParkA", "Reddit_HealthTipsCn", "Reddit_HoroScopeCn", "Reddit_HotGamesCn", "Reddit_HotGamesEn", "Reddit_Joke", "Reddit_JokeA", "Reddit_JokeBoxEn", "Reddit_JokeBoxKidsEn", "Reddit_JokeEn", "Reddit_JokesBoxA", "Reddit_LY_EnglishQuotesCn", "Reddit_LY_JokesEn", "Reddit_LY_QQStatus", "Reddit_LYqqstatusA", "Reddit_LeStoryCNA", "Reddit_LifeTipsCn", "Reddit_LittleJohnnyJokes", "Reddit_LoveQuotesCn", "Reddit_MyPhotoCNA", "Reddit_OneDirection", "Reddit_QuestionEn", "Reddit_QuotesCn", "Reddit_RateBoxEn", "Reddit_RateBoxEnA", "Reddit_RateMeA", "Reddit_SY_HoroScopeEn", "Reddit_SexFactsEn", "Reddit_SexfactsA", "Reddit_ShareMVCNA", "Reddit_SweetMoment", "Reddit_SweetMomentA", "Reddit_WeiQingShu", "Reddit_WeiXiaoShuo", "Reddit_WeiqingshuA", "Reddit_WeixiaoshuoA", "Reddit_WeixiaotuCNA", "Reddit_dirtyjokesenA", "Reddit_funnypicturesCNA", "Reddit_funnypicturesENA", "Reddit_horoscope_CNA", "Reddit_horoscope_ENA", "Reddit_needadvicecnA", "Reddit_quotescnA", "Ringtone_Builder", "RingtonesBox", "SMSViewer", "SoundEffectsAllIn1", "SpeedU", "TM_Themes", "WM_FantasyWallpapers", "WM_iOS7Wallpapers", "iDownload", "sns_photoeditor"]

#APPS = ["1P_1080Wallpapers", "1P_47inchWallpapers", "1P_BigWallpapers", "1P_FineWallpapers", "1P_GlassWallpapers", "1P_GrassWallpapers", "1P_HugeWallpapers", "1P_LargeWallpapers", "1P_LockWallpapers", "1P_PanoramaWallpaperWallpapers", "1P_PhabletWallpapers", "1P_RainbowWallpapers", "1P_SapphireWallpapers", "1P_SmartWallpapers", "1P_ThemeWallpapers", "1P_Wallpapers", "1P_Wallpapers8", "1P_WallpapersHD", "1P_iOS8Wallpapers", "1P_iPhone6Wallpapers", "4P_CoolColorWallpapers", "4P_FunnyPix", "4P_FunnyPix_iPad", "4P_HarryPotterWallpapers", "4P_HelloKittyWallpapers", "4P_JustinBieberWallpapers", "4P_LockScreenWallpapers", "4P_LoveQuotesWallpapers", "9P_iOS7OnlyWallpaper_iPad", "9P_iOS7Wallpapershd", "9P_iOS7Wallpapershd_iPad", "9P_iPhone5Wallpapers", "9P_iPhone5Wallpapers_iPad", "AR_AppRank", "AR_AppRed", "AR_COC", "AR_Common", "AR_CoverFree", "AR_CoverFreeCN", "AR_Desktop", "AR_Explore", "AR_FreeAppToday", "AR_GameFreeCN", "AR_GameLimitedFree", "AR_GameVideo", "AR_HDFreeCN", "AR_IconFreeCN", "AR_IconFreeCN_3_7", "AR_IconFreeCN_3_8", "AR_IconSortLimitedFree", "AR_LOL", "AR_Launcher", "AR_LimitedFree", "AR_ListSortLimitedFree", "AR_MY", "AR_MustAppsCN", "AR_SpeedU1", "AR_WePlay", "AR_XMStatusCN", "AR_XYPlus", "AR_XianYou", "AR_XianYouAnswer", "AR_XianYouAnswerHD", "AR_XianYouApp", "AR_XianYouAsk", "AR_XianYouGame", "AR_XianYouGameNews", "AR_XianYouVideo", "AT_IconER", "AT_IconERLite", "AT_IconERPro", "AV_AutoVideo", "AV_BeautyGirl", "AV_FunnyTime", "AY_Ipomelo", "AY_PhoneCase", "AY_Shoujike", "ChristmasRingtones", "Downloader", "FBPhoto", "FunnyRingtones", "RT_360Ringtones", "RT_365Ringtones", "RT_EverRingtones", "RT_KuaiBo", "RT_PerfectRingtones", "RT_PowerRingtones", "RT_RingtonesBox", "RT_RingtonesBoxPro", "RT_RingtonesBox_Android", "RT_RingtonesStyle", "RT_iOS8Ringtones", "RT_iPhone6PlusRingtones", "RT_iPhone6Ringtones"]

def import_data():
    global HOST, json
    s = requests.Session()

    #init user
    r = s.post("%s/op/user/init"%(HOST), json={"name": "admin", "pass_code":"111111", "aux_info":"{}"})
    if r.json()["status"]==False:
        #print r.json()
        pass

    #login
    r = s.post("%s/op/login"%(HOST), json={"name": "admin", "pass_code":"111111"})
    if r.json()["status"]==False:
        print r.json()
        
    #add user
    aux_info = {}
    for user in ["rahuahua", "ldmiao", "leo", "lishaohua", "gaoqingzu", "zhaoyong", "zhanhao", "kdr2"]:
        for i in range(30):
            if i==0:
                username = user
            else:
                username = "%s%02d"%(user, i)

            json_data = {"name": username, "pass_code":"111111", "aux_info":json.dumps(aux_info)}
            r = s.post("%s/op/user"%(HOST), json=json_data)
            if r.json()["status"]==False:
                print r.json()

    #add app
    aux_info = {
                    "link"  : "https://itunes.apple.com/cn/app/chao-gao-qing-bi-zhi-zui-re/id384922950?mt=8",
                    "title" : "超高清壁纸",
                    "icon"  : "http://a1.mzstatic.com/us/r30/Purple5/v4/39/0f/23/390f236b-1228-b1a7-3315-a456115b4381/icon175x175.png"
               }

    global APPS
    for appname in APPS:
        json_data = {"name": appname, "type":"real", "aux_info":json.dumps(aux_info)}
        r = s.post("%s/op/app"%(HOST), json=json_data)
        if r.json()["status"]==False:
            print r.json()

    #get appname -> app_key
    app_keys = {}
    r = s.get("%s/op/apps/all/1/10000"%(HOST))
    d = r.json()
    if d["status"]==True:
        for app in d["data"]["list"]:
            app_keys[app["name"]] = app["key"]

    print app_keys

    #add config
    for app_name, app_key in app_keys.items():
        print app_name, app_key
        r = requests.get("http://conf2.appwill.com/conf?app=%s"%(app_name))
        d = r.json()
        if d:
            for k,v in d.items():
                v_type = "string"
                if isinstance(v, int):
                    v_type = "int"
                    v = "%d"%v
                json = {"k": k, "v":v, "v_type":v_type, "app_key":app_key, "desc":"", "status":1}
                r = s.post("%s/op/config"%(HOST), json=json)
                if r.json()["status"]==False:
                    print k, v, v_type, r.json()

def check_data():
    global APPS, HOST, json
    for app_name in APPS:
        r1 = requests.get("%s/conf?app=%s"%(HOST, app_name))
        d1 = r1.json()
        
        r2 = requests.get("http://conf1.appwill.com/conf?app=%s"%(app_name))
        d2 = r2.json()
        
        r3 = requests.get("http://conf2.appwill.com/conf?app=%s"%(app_name))
        d3 = r3.json()

        for k,v1 in d1.items():
            v2 = d2.get(k, "instafig_default_value")
            if v2!=v1:
                print app_name, k, v1, v2

        for k,v2 in d2.items():
            v1 = d1.get(k, "instafig_default_value")
            if v1!=v2:
                print app_name, k, v1, v2

        for k,v1 in d1.items():
            v3 = d3.get(k, "instafig_default_value")
            if v1!=v3:
                print app_name, k, v1, v3

        for k,v3 in d3.items():
            v1 = d1.get(k, "instafig_default_value")
            if v1!=v3:
                print app_name, k, v1, v3

if __name__=='__main__':
    import_data()
    #check_data()
