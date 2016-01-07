var infoSysDirectives = angular.module('infoSysDirectives', []);

infoSysDirectives.directive('informer', function () {
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

infoSysDirectives.directive('leftpanel', function () {
    return {
        controller: function($scope, $location, $browser) {
            $scope.items = [
                //{text : "Объекты", path: "objects"},
                //{text : "Редактор объектов", path: "objects-edit"},
                //{text : "Конфигурация", path: "configs"},
                //{text : "Сотрудники", path: "accounts"},
                //{text : "Отделы", path: "groups"},
                {text : "База данных", path: "db"}
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

infoSysDirectives.directive('button', function () {
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

infoSysDirectives.directive('propcontrol', function () {
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

infoSysDirectives.directive("dropdown", function($rootScope) {
    return {
        restrict: "E",
        templateUrl: "dropdown-control",
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
                return item[$scope.property] === $scope.selected[$scope.property];
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
                $scope.isPlaceholder = $scope.selected[$scope.property] === undefined;
                $scope.display = $scope.selected[$scope.property];
            });
        }
    }
});