
var conf_condition;//配置项中点击的value
var cond_data = new Array();
var show_condition;//ifcond展示
//条件-值 列表
var ShowConditions = function (resolve, reject) {
    var template_url = 'components/show_conditions.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        show_condition = 0;
                        if (conf_condition != null) {
                            try {
                                cond_data = [];
                                if (conf_condition.length != undefined) {
                                    for (i in conf_condition) {
                                        cond_data.push(conf_condition[i]);
                                    }
                                } else {
                                    cond_data = [];
                                    cond_data.push(conf_condition)
                                }

                            } catch (e) {

                            }
                        }
                        console.log(cond_data)
                        return {
                            data_condition: cond_data,
                            condition_data: null,
                            logic_funcs: LOGIC_FUNCS,      //逻辑操作func列表
                            leaf_func_titles: LEAF_FUNC_TITLES, //func对应的显示名
                            symbols: SYMBOLS,          //Symbol列表
                            symbol_titles: SYMBOL_TITLES,    //Symbol对应的显示名  //Symbol对应的可用func列表
                            symbol_funcs: SYMBOL_FUNCS,    //Symbol对应的可用func列表
                            readonly: 'true',
                            cond_num: '',
                            index: -1,
                            num: ''
                        };
                    }
                },
                computed: {
                    //判断单层和多层
                    is_leaf_cond: function () {
                        var func = this.condition_data.func;
                        return LEAF_FUNC_TITLES[func] != undefined;
                    },
                    //判断能否添加条件
                    can_add_sub_cond: function () {
                        var func = this.condition_data.func;
                        return func == 'and' || func == 'or';
                    },
                    height_for_func: function () {
                        return 6 + 20 * (this.condition_data.arguments.length - 2);
                    }
                },
                methods: {
                    //增加条件
                    add_config: function () {
                        cond_data.push({ "condition": { "arguments": [{ "Symbol": "LANG" }, "en"], "func": "str!=?" }, "value": 0 });
                    },
                    //删除条件
                    delete_config: function (idx) {
                        cond_data.splice(idx, 1);
                    },
                    //保存条件修改
                    save_sub_condition: function (condition) {
                        cond_data[this.num].condition = condition;
                        $('#myModal2').modal('hide');
                    },
                    add_sub_cond: function (func) {   
                        //将增加的条件添加到condition数组
                        //存在多层时  
                        if (["and", "or"].indexOf(func) != -1) {
                            this.condition_data.arguments.push({ "arguments": [{ "Symbol": "APP_VERSION" }, "1.0"], "func": "ver=" });
                        } else if (func == 'not') {
                            this.condition_data.func = 'and';
                            this.condition_data.arguments.push({ "arguments": [{ "Symbol": "APP_VERSION" }, "1.0"], "func": "ver=" });
                        } else {
                            //单层时
                            this.condition_data = { "arguments": [this.condition_data, { "arguments": [{ "Symbol": "APP_VERSION" }, "1.0"], "func": "ver=" }], "func": "and" };
                        }
                    },
                    del_sub_cond: function (len, idx) {
                        //判断idx是否是初始条件
                        if (idx < this.cond_num) {
                            this.cond_num = this.cond_num - 1;
                        }
                        if (len > 2) {
                            this.condition_data.arguments.splice(idx, 1);
                        } else if (len == 2) {
                            this.condition_data.arguments.splice(idx, 1);
                            if (this.condition_data.func == 'not') {
                                this.condition_data.func = 'not';
                            } else {
                                this.condition_data = this.condition_data.arguments[0];
                            }
                        }
                    },
                    change_func: function () {
                        var func = this.condition_data.func;
                        if (func == "and") {
                            //当为and时，点击变为or
                            this.condition_data.func = "or";
                        } else {
                            //当为or时，点击变为and
                            this.condition_data.func = "and";
                        }
                    },
                    //编辑条件
                    show_condition_edit: function (idx, conditions) {
                        this.num = idx;
                        //判断是否是单层,计算初始条件数
                        if (LEAF_FUNC_TITLES[conditions.condition.func] != undefined) {
                            this.cond_num = 1;
                        } else {
                            this.cond_num = conditions.condition.arguments.length;
                        }
                        conf_condition = conditions;
                        this.condition_data = clone(conditions.condition)
                    },
                    //保存
                    save_condition: function () {
                        for (i in this.data_condition) {
                            this.data_condition[i].value = this.data_condition[i].value
                        }
                        var conf_value = JSON.stringify(JSON.stringify(cond_data));
                        fetch('/op/config', {
                            method: 'PUT',
                            body: JSON.stringify({
                                "key": conf_data_array[2],
                                "k": conf_data_array[5],
                                "v": conf_value,
                                "v_type": conf_data_array[1],
                                "app_key": conf_data_array[3],
                                "des": conf_data_array[4],
                                "status": parseInt(conf_data_array[6])
                            }),
                            credentials: 'same-origin'
                        }).then(function (response) {
                            return response.json();
                        }).then(function (data) {
                            if (data.status == true) {
                                $('#myModal1').modal('hide');
                                router.go("/apps");
                            }
                        })
                    }
                }
            });
        });
};

var conf_data_array = new Array();//配置项各项
//配置列表
var AppConfig = function (resolve, reject) {
    var template_url = 'components/app_config_list.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        show_condition = 1;
                        var condition = new Array();
                        var app_key = transition.to.params.app_key;
                        var app_name = transition.to.params.app_name;
                        var data_url = '/op/configs/' + app_key;
                        return fetch(data_url, {
                            methods: 'GET',
                            credentials: 'same-origin'
                        }).then(function (data_resp) {
                            return data_resp.json();
                        }).then(function (data) {
                            var exit = new Array();//给条件标识，0代表value存在多个条件，1代表单个条件(可直接编辑)                                                       
                            for (i in data.data) {
                                var config_v = data.data[i];
                                try {
                                    config_v.v = JSON.parse(config_v.v)
                                    if (JSON.parse(config_v.v).length != undefined) {
                                        exit.push(0)//数组条件
                                    } else {
                                        if (JSON.parse(config_v.v).value != undefined) {
                                            exit.push(0);//多个条件
                                        } else {
                                            exit.push(1)
                                        }
                                    }
                                } catch (e) {
                                    config_v.v = config_v.v
                                    exit.push(1);//单个条件
                                }
                                try {
                                    if (JSON.parse(config_v.v).length == undefined) {
                                        if (JSON.parse(config_v.v).value != undefined) {
                                            condition.push([JSON.parse(config_v.v)]);
                                        } else {
                                            condition.push(0);
                                        }
                                    } else {
                                        condition.push(JSON.parse(config_v.v));
                                    }
                                } catch (e) {
                                    condition.push(0);
                                }
                            }
                            return {
                                app_key: app_key,
                                app_name: app_name,
                                data: data,
                                index: '-1',
                                conf_value_key: '',
                                conf_value_status: '',
                                conf_value_des: '',
                                exit_condition: exit,
                                cond_data: condition,
                                key_error: '',
                                value_type_error: ''
                            };
                        })
                    }
                },
                data: {
                    conf_key: '',
                    conf_value: '',
                    conf_value_type: '',
                    conf_des: '',
                },
                methods: {
                    //单个条件的value
                    edit_value: function (idx) {
                        this.index = idx;
                    },
                    //单个条件的value编辑
                    modify_value: function (idx) {
                        fetch('/op/config', {
                            method: 'PUT',
                            body: JSON.stringify({
                                "key": this.data.data[idx].key,
                                "k": this.data.data[idx].k,
                                "v": this.data.data[idx].v,
                                "v_type": this.data.data[idx].v_type,
                                "app_key": this.data.data[idx].app_key,
                                "des": this.data.data[idx].des,
                                "status": this.data.data[idx].status
                            }),
                            credentials: 'same-origin'
                        }).then(function (response) {
                            return response.json();
                        }).then(function (data) {
                        })
                    },
                    //多个条件时value编辑
                    modify_config: function (idx) {
                        console.log(idx)
                        conf_condition = JSON.parse(this.data.data[idx].v);
                        conf_data_array[5] = this.data.data[idx].k;
                        conf_data_array[4] = this.data.data[idx].des;
                        conf_data_array[1] = this.data.data[idx].v_type;
                        conf_data_array[2] = this.data.data[idx].key;
                        conf_data_array[3] = this.data.data[idx].app_key;
                        conf_data_array[6] = this.data.data[idx].status;
                        router.go('/show_conditions')
                    },
                    //新增配置
                    add_config: function () {
                        var vm = this;
                        var conf_value;
                        if (JSON.parse(this.conf_value).value != undefined) {
                            conf_value = JSON.stringify(this.conf_value);
                        } else {
                            conf_value = this.conf_value;
                        }
                        fetch('/op/config', {
                            method: 'POST',
                            body: JSON.stringify({
                                "k": this.conf_key,
                                "v": conf_value,
                                "v_type": this.conf_value_type,
                                "app_key": this.app_key,
                                "des": this.conf_des
                            }),
                            credentials: 'same-origin'
                        }).then(function (response) {
                            return response.json();
                        }).then(function (data) {
                            console.log(data.status)
                            if (data.status == true) {
                                window.location.reload();
                                $('#myModal').modal('hide');
                            } else if (is.startWith(data.msg, 'config key has existed')) {
                                vm.key_error = "key已存在";
                                vm.value_type_error = "";
                            } else {
                                vm.key_error = "";
                                vm.value_type_error = "value类型错误";
                            }
                        });
                    },
                    //修改配置
                    modify_condition: function (idx) {
                        this.conf_value_key = this.data.data[idx].k;
                        this.conf_value_des = this.data.data[idx].des;
                        this.conf_value_status = this.data.data[idx].status;
                        conf_data_array[0] = this.data.data[idx].v;
                        conf_data_array[1] = this.data.data[idx].v_type;
                        conf_data_array[2] = this.data.data[idx].key;
                        conf_data_array[3] = this.data.data[idx].app_key;
                    },
                    //保存配置的修改
                    sava_modify: function () {
                        fetch('/op/config', {
                            method: 'PUT',
                            body: JSON.stringify({
                                "key": conf_data_array[2],
                                "k": this.conf_value_key,
                                "v": conf_data_array[0],
                                "v_type": conf_data_array[1],
                                "app_key": conf_data_array[3],
                                "des": this.conf_value_des,
                                "status": parseInt(this.conf_value_status)
                            }),
                            credentials: 'same-origin'
                        }).then(function (response) {
                            return response.json();
                        }).then(function (data) {
                            if (data.status == true) {                               
                                $('#myModal1').modal('hide');
                                window.location.reload();
                            }
                        })
                    },
                    modify_history: function () {
                        router.go('');
                    }
                }
            })
        })
}

//Vue component - ifcond
Vue.component('ifcond', function (resolve, reject) {
    var template_url = 'components/ifcond.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                props: ['model'],
                data: function () {
                    return {
                        logic_funcs: LOGIC_FUNCS,      //逻辑操作func列表
                        leaf_func_titles: LEAF_FUNC_TITLES, //func对应的显示名
                        symbols: SYMBOLS,          //Symbol列表
                        symbol_titles: SYMBOL_TITLES,    //Symbol对应的显示名  //Symbol对应的可用func列表
                        symbol_funcs: SYMBOL_FUNCS,    //Symbol对应的可用func列表
                        show_name: show_condition
                    }
                },
                computed: {
                    is_leaf_cond: function () {
                        var func = this.model.func;
                        return LEAF_FUNC_TITLES[func] != undefined;
                    },
                    fr: function () {
                        var func = this.model.func;
                        var args = this.model.arguments;
                        var ret = null;
                        if (LEAF_FUNC_TITLES[func] != undefined) {
                            ret = [args[0], { operator: func }, { value: args[1] }];
                        } else if (func == "not") {
                            ret = [{ operator: func }, args[0]];
                        } else if (["and", "or"].indexOf(func) != -1) {
                            ret = [];
                            for (var i = 0; i < args.length; i++) {
                                ret.push(args[i]);
                                if (i != args.length - 1) {
                                    ret.push({ operator: func });
                                }
                            }
                        }
                        return ret;
                    }
                }
            });
        });
});

//Vue component - cond
Vue.component('cond', function (resolve, reject) {
    var template_url = 'components/cond.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                props: ['model'],
                computed: {
                    subConds: function () {
                    }
                }
            });
        });
});