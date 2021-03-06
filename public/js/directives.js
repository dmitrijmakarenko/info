var accessSettingsDirectives = angular.module('accessSettingsDirectives', []);

accessSettingsDirectives.directive('informer', function () {
    return {
        controller: function($scope, informerControl, $timeout) {
            $scope.show = false;
            $scope.informerControl = informerControl;

            $scope.$watch('informerControl.status', toggledisplay);
            $scope.$watch('informerControl.message', toggledisplay);
            $scope.$watch('show', hide);

            function toggledisplay() {
                $scope.show = !!($scope.informerControl.status && $scope.informerControl.message);
            }
            function hide(value) {
                if (value) {
                    $timeout(function() {
                        $scope.show = false;
                        $scope.informerControl.clear();
                    }, 3000);
                }
            }
        },
        template:
        '<div class="informer {{informerControl.status}}" ng-show="show">' +
            '<div class="msg">{{informerControl.message}}</div>' +
        '</div>'
    }
});

accessSettingsDirectives.directive('spinner', function () {
    return {
        controller: function($scope, spinnerControl) {
            $scope.spinnerControlShow = false;
            $scope.spinnerControl = spinnerControl;

            $scope.$watch('spinnerControl.spinnerShow', toggle);
            function toggle() {
                $scope.spinnerControlShow = spinnerControl.spinnerShow;
            }
        },
        template:
        '<div class="spinner-loading" ng-show="spinnerControlShow">' +
            '<div class="shade"></div>' +
            '<img src="/public/img/loading.gif" alt="Загрузка..." class="spinner"/>' +
        '</div>'
    }
});

accessSettingsDirectives.directive('leftpanel', function () {
    return {
        controller: function($scope, $location, $browser) {
            $scope.items = [
                {text : "Таблицы", path: "tables"},
                {text : "Пользователи", path: "accounts"},
                {text : "Группы", path: "groups"},
                {text : "Метки", path: "rules"},
                //{text : "База данных", path: "db"},
                //{text : "VCS", path: "vcs"},
                {text : "Тест", path: "test"}
            ];
            $scope.go = function(val) {
                window.location = "#/" + val.path;
                for (var i = 0; i < $scope.items.length; i++) {
                    if ($scope.items[i].path == val.path) {
                        $scope.items[i].clazz = "active";
                    } else {
                        $scope.items[i].clazz = null;
                    }
                }

            };
            $browser.onUrlChange(function (newUrl) {
                var getPath = /(#\/.*)/,
                    match = newUrl.match(getPath),
                    path;
                if (match.length > 1) {
                    path = match[0];
                    path = path.replace(/#\//g, "");
                    changePath(path);
                }
            });

            var path = $location.$$path;
            path = path.replace(/\//g, "");
            changePath(path);

            function changePath(path) {
                var getRoot = /(.*\/)/,
                    match = path.match(getRoot);
                if (match && match.length > 1) {
                    path = match[0];
                    path = path.replace(/\//g, "");
                }
                for (var i = 0; i < $scope.items.length; i++) {
                    if ($scope.items[i].path == path) {
                        $scope.items[i].clazz = "active";
                    } else {
                        $scope.items[i].clazz = null;
                    }
                }
            }
        },
        template:
        '<div id="left-panel">' +
            '<div class="logo"></div>' +
            '<div ng-repeat="item in items" ng-click="go(item)" class="menu-item {{item.clazz}}">' +
                '<div class="text">{{item.text}}</div>' +
            '</div>' +
        '</div>'
    }
});

accessSettingsDirectives.directive('button', function () {
    return {
        scope: true,
        controller: function($scope) {
            $scope.onClick = function() {
                if ($scope.func && $scope.$parent[$scope.func]) {
                    if ($scope.loading) {
                        $scope.btnLoading = true;
                        $scope.$parent[$scope.func].apply(null, $scope.params);
                        $scope.btnLoading = false;
                    } else {
                        $scope.$parent[$scope.func].apply(null, $scope.params);
                    }
                }
            };
        },
        link: function($scope, element, attrs) {
            if (attrs.text) {
                $scope.text = attrs.text;
            }
            if (attrs.clazz) {
                $scope.clazz = attrs.clazz;
            }
            if (attrs.func) {
                $scope.func = attrs.func;
            }
            if (attrs.params) {
                $scope.params = attrs.params.split(',');
            } else {
                $scope.params = [];
            }
            $scope.loading = (attrs.loading ? true : false);
        },
        template: '<div class="button {{clazz}}" ng-click="onClick()">' +
        '<div class="btn-spinner-wrap" ng-show="btnLoading"><div class="btn-spinner"></div></div>' +
        '<div class="btn-label" ng-show="!btnLoading">{{text}}</div>' +
        '</div>'
    }
});

accessSettingsDirectives.directive('propcontrol', function () {
    return {
        controller: function($scope, ngDialog) {
            $scope.types = [
                {value: "int", text: "integer"},
                {value: "char", text: "char"},
                {value: "double", text: "double"}
            ];
            $scope.typeSelected = {};


            $scope.addNew = function() {
                ngDialog.open({
                    template: 'edit-prop-dialog',
                    controller: 'editPropCntl',
                    disableAnimation: true,
                    showClose: false,
                    scope: $scope
                });
            };

            $scope.edit = function(id) {
                var params = {};
                for (var i = 0; i < $scope.properties.length; i++) {
                    if ($scope.properties[i].id == id) {
                        params.id = id;
                        var type = {};
                        for (var j = 0; j < $scope.types.length; j++) {
                            if ($scope.types[j].value == $scope.properties[i].type) {
                                type = $scope.types[j];
                            }
                        }
                        params.type = type;
                        params.desc = $scope.properties[i].desc||"";
                    }
                }
                ngDialog.open({
                    template: 'edit-prop-dialog',
                    controller: 'editPropCntl',
                    disableAnimation: true,
                    showClose: false,
                    data: {mode: "edit", params: params},
                    scope: $scope
                });
            };

            $scope.remove = function(id) {
                var tmp = [];
                for (var i = 0; i < $scope.properties.length; i++) {
                    if ($scope.properties[i].id != id) {
                        tmp.push($scope.properties[i]);
                    }
                }
                $scope.properties = tmp;
            };
        },
        templateUrl: 'property-control'
    }
});

accessSettingsDirectives.directive("dropdown", function($rootScope) {
    return {
        restrict: "E",
        scope: {
            placeholder: "@",
            list: "=",
            selected: "=",
            property: "@"
        },
        link: function($scope) {
            $scope.listVisible = false;
            $scope.isPlaceholder = true;

            $scope.select = function(item) {
                $scope.isPlaceholder = false;
                $scope.selected = item;
                if ($scope.onChange !== undefined)
                    $scope.onChange(item);
            };

            $scope.isSelected = function(item) {
                if ($scope.selected) {
                    return item[$scope.property] === $scope.selected[$scope.property];
                } else {
                    return false
                }
            };

            $scope.show = function() {
                $scope.listVisible = true;
            };

            $rootScope.$on("documentClicked", function(inner, target) {
                if (!$(target[0]).is(".dropdown-display.clicked") && !$(target[0]).parents(".dropdown-display.clicked").length > 0)
                    $scope.$apply(function() {
                        $scope.listVisible = false;
                    });
            });

            $scope.$watch("selected", function() {
                if ($scope.selected) {
                    $scope.isPlaceholder = $scope.selected[$scope.property] === undefined;
                    $scope.display = $scope.selected[$scope.property];
                }
            });
        },
        template:
        '<div class="dropdown-container" ng-class="{ show: listVisible }">' +
            '<div class="dropdown-display" ng-click="show();" ng-class="{ clicked: listVisible }">' +
                '<span ng-if="!isPlaceholder">{{display}}</span>' +
                '<span class="placeholder" ng-if="isPlaceholder">{{placeholder}}</span>' +
                '<i class="fa fa-angle-down"></i>' +
            '</div>' +
            '<div class="dropdown-list">' +
                '<div>' +
                    '<div ng-repeat="item in list" ng-click="select(item)" ng-class="{ selected: isSelected(item) }">' +
                        '<span>{{property !== undefined ? item[property] : item}}</span>' +
                        '<i class="fa fa-check"></i>' +
                    '</div>' +
                '</div>' +
            '</div>' +
        '</div>'
    }
});