//全局配置
var GlobalConf = function (resolve, reject) {
    var template_url = 'components/global_conf.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        return fetch('/op/webhooks/global', { credentials: 'same-origin' })
                            .then(function (data_resp) {
                                return {
                                    data: data_resp.json()
                                };
                            }).then(function (data) {
                            });
                    }
                },
                data: {
                },
                methods: {
                    //增加Webhook
                    add_webhook: function () {
                    }
                }
            });
        });
};

//服务状态
var ServerStatus = function (resolve, reject) {
    var template_url = 'components/server_status.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        return fetch('/op/nodes', { credentials: 'same-origin' })
                            .then(function (data_resp) {
                                return data_resp.json();
                            }).then(function (data) {
                                return {
                                    data: data
                                };
                            })
                    }
                }
            });
        });
};