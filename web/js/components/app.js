//应用列表
var Apps = function (resolve, reject) {
    var template_url = 'components/apps.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        var page = transition.to.params.page;
                        var path;
                        if (is.startWith(transition.to.path, '/apps/all')) {
                            path = '/op/apps/all/' + transition.to.params.page + '/' + COUNT_PER_PAGE;
                        } else if (is.startWith(transition.to.path, '/apps/search')) {
                            path = '/op/apps/search?q=' + transition.to.params.keyword;
                        } else if (is.startWith(transition.to.path, '/apps/user')) {
                            path = '/op/apps/user/' + transition.to.params.user_key;
                        }
                        return fetch(path, { method: 'GET', credentials: 'same-origin' })
                            .then(function (data_resp) {
                                return data_resp.json();
                            }).then(function (data) {
                                if (data && data.data != null) {
                                    var app;
                                    var page_count = null;
                                    var app_data;
                                    if (path == '/op/apps/all/' + page + '/' + COUNT_PER_PAGE) {
                                        for (i in data.data.list) {
                                            app_data = data.data.list;
                                            page_count = 0;
                                            app = data.data.list[i];
                                            try {
                                                app.aux_info = JSON.parse(app.aux_info)
                                            } catch (e) {
                                                app.aux_info = {};
                                            }
                                            if (data.data.total_count % COUNT_PER_PAGE == 0) {
                                                page_count = data.data.total_count / COUNT_PER_PAGE;
                                            } else {
                                                page_count = Math.floor(data.data.total_count / COUNT_PER_PAGE) + 1;
                                            }
                                        }
                                    } else {
                                        for (i in data.data) {
                                            app_data = data.data;
                                            app = data.data[i];
                                            try {
                                                app.aux_info = JSON.parse(app.aux_info)
                                            } catch (e) {
                                                app.aux_info = {};
                                            }
                                        }
                                    }
                                }
                                return {
                                    data: app_data,
                                    app: null,
                                    operater: '',
                                    modal_title: '',
                                    judge_value: 1,
                                    page_count: page_count,
                                    err_msg: null,
                                    keyword: null,
                                }
                            })
                    }
                },
                methods: {
                    //保存App数据
                    save_data: function () {
                        var vm = this;
                        if (vm.judge_value == 1) {
                            //增加APP
                            fetch('/op/app', {
                                method: 'POST',
                                body: JSON.stringify({
                                    "name": vm.app.name,
                                    "cloudname": vm.app.cloudname,
                                    "domain":vm.app.domain,
                                    "bucket":vm.app.bucket,
                                    "dateformat":vm.app.dateformat,
                                    "aux_info": JSON.stringify(vm.app.aux_info)
                                }),
                                credentials: 'same-origin'
                            }).then(function (response) {
                                return response.json();
                            }).then(function (data) {
                                if (data.status == true) {
                                    $('#myModal').modal('hide');
                                } else {
                                    if (is.startWith(data.msg, 'appname already exists')) {
                                        vm.err_msg = "应用名已存在，请重新输入。";
                                    }
                                }
                            });
                        } else if (vm.judge_value == 2) {
                            //修改APP
                            fetch('/op/app', {
                                method: 'PUT',
                                body: JSON.stringify({
                                    "key": vm.app.key,
                                    "name": vm.app.name,
                                    "cloudname": vm.app.cloudname,
                                    "domain":vm.app.domain,
                                    "bucket":vm.app.bucket,
                                    "dateformat": vm.app.dateformat,
                                    "aux_info": JSON.stringify(vm.app.aux_info)
                                }),
                                credentials: 'same-origin'
                            }).then(function (response) {
                                return response.json();
                            }).then(function (data) {
                                if (data.status == true) {
                                    $('#myModal').modal('hide');
                                } else {
                                    if (is.startWith(data.msg, 'appname already exists')) {
                                        vm.err_msg = "应用名已存在，请重新输入。";
                                    }
                                }
                            });
                        }
                    },
                    //新增App
                    show_add_app: function () {
                        this.err_msg = null;
                        this.judge_value = 1;
                        this.app = {
                            "name": "",
                            "type": "real",
                            "aux_info": {
                                "title": "",
                                "icon": "",
                                "link": ""
                            }
                        };
                        this.modal_title = '增加应用';
                        this.operater = '增加';
                    },
                    //修改App
                    show_edit_app: function (app) {
                        this.err_msg = null;
                        this.judge_value = 2;
                        this.app = clone(app);
                        this.modal_title = '修改应用';
                        this.operater = '保存';
                    },
                    //克隆App
                    show_clone_app: function (app) {
                        this.err_msg = null;
                        this.judge_value = 1;
                        this.app = clone(app);
                        this.app.key = null;
                        this.modal_title = '克隆应用';
                        this.operater = '克隆';
                    },
                    all_apps: function () {
                        router.go('/apps')
                    },
                    search_app: function (keyword) {
                        router.go('apps/search/' + keyword);
                    }
                }
            });
        });
};