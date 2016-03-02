accessSettings.controller('entitiesCntl', function ($scope, Entity) {

    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getObject = function(entity) {
        window.location = "#/objects-edit/" + entity;
    }
});

accessSettings.controller('entityCntl', function ($scope, $routeParams, Entity, ngDialog) {
    $scope.loading = false;
    var entityId = $routeParams.objectId;

    $scope.sqlActions = [
        {action: "SELECT", service: "GetSqlSelect", text: "Просмотр данных", show: false},
        {action: "UPDATE", service: "GetSqlUpdate", text: "Редактирование данных", show: false},
        {action: "INSERT", service: "GetSqlInsert", text: "Добавление данных", show: false},
        {action: "DELETE", service: "GetSqlDelete", text: "Удаление данных", show: false}
    ];

    $scope.toggle = function(index) {
        $scope.sqlActions[index].show = !$scope.sqlActions[index].show;
    };
    $scope.getQuery = function(index) {
        var params = {};
        params.table = entityId;
        Entity[$scope.sqlActions[index].service].go(params, function(data) {
            $scope.sqlActions[index].query = data.query;
        });
        /*if (action == "SELECT") {
            Entity.GetSqlSelect.go({table: entityId}, function(data) {
                $scope.sqlActions[index].query = data.query;
            });
        } else if (action == "UPDATE") {
            Entity.GetSqlUpdate.go({table: entityId}, function(data) {
                $scope.sqlActions[index].query = data.query;
            });
        } else if (action == "DELETE") {
            Entity.GetSqlDelete.go({table: entityId}, function(data) {
                $scope.sqlActions[index].query = data.query;
            });
        } else if (action == "INSERT") {
            Entity.GetSqlInsert.go({table: entityId}, function(data) {
                $scope.sqlActions[index].query = data.query;
            });
        }*/
    };

    if (entityId != "!new") {
        $scope.loading = true;
        Entity.Get.go({id: entityId}, function(data) {
            $scope.loading = false;
            $scope.id = data.id;
            $scope.name = data.name;
            $scope.properties = data.props||[];
            $scope.createMode = true;
        });
    } else {
        $scope.id = "";
        $scope.name = "";
        $scope.properties = [];
        $scope.createMode = false;
    }

    $scope.saveSettings = function() {
        if (entityId == "!new") {
            if (validate()) {
                Entity.Create.go({id: $scope.id, name: $scope.name, props: JSON.stringify($scope.properties)}, function(data) {
                    if (data.error) {
                        $scope.showErrorMsg(data.error);
                    } else {
                        $scope.showSuccessMsg("Сохранено");
                        window.location = "#/entities";
                    }
                });
            }
        } else {
            if (validate()) {
                Entity.Update.go({id: $scope.id, name: $scope.name, props: JSON.stringify($scope.properties)}, function(data) {
                    if (data.error) {
                        $scope.showErrorMsg(data.error);
                    } else {
                        $scope.showSuccessMsg("Сохранено");
                        window.location = "#/entities";
                    }
                });
            }
        }
    };

    $scope.entityRemove = function() {
        ngDialog.openConfirm({
            template: 'confirm-delete-dialog',
            disableAnimation: true,
            showClose: false
        }).then(function () {
            Entity.Remove.go({id: entityId}, function(data) {
                if (data.error) {
                    $scope.showErrorMsg(data.error);
                } else {
                    $scope.showSuccessMsg("Удалено");
                    window.location = "#/entities";
                }
            });
        }, function () {

        });
    };

    function validate() {
        var valid = true;
        if ($scope.id == "") {
            valid = false;
        }

        if ($scope.name == "") {
            valid = false;
        }

        return valid;
    }

});

accessSettings.controller('objectsCntl', function ($scope, Entity) {
    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getObject = function(entity) {
        window.location = "#/objects/" + entity;
    }
});

accessSettings.controller('objectCntl', function ($scope, $routeParams, Entity) {
    $scope.loading = true;
    var entityId = $routeParams.objectId;

    Entity.Get.go({id: entityId, returnData: true}, function(data) {
        if (!data.error) {
            $scope.entity = data;
            if (!$scope.entity.rows) $scope.entity.rows = [];
        }
        console.log(data);
    });
});

accessSettings.controller('editPropCntl', function ($scope, Entity, ngDialog) {
    var editMode = false,
        prevId;
    if ($scope.ngDialogData && $scope.ngDialogData.mode && $scope.ngDialogData.mode == "edit") editMode = true;

    if (editMode) {
        $scope.id = $scope.ngDialogData.params.id||"";
        $scope.typeSelected = $scope.ngDialogData.params.type||{};
        $scope.desc = $scope.ngDialogData.params.desc||"";
        prevId = $scope.id;
    } else {
        $scope.id = "";
        $scope.type = "";
        $scope.desc = "";
    }
    $scope.saveProperty = function() {
        if ($scope.id != "" && $scope.typeSelected) {
            if (editMode) {
                var tmp = $scope.properties.slice(0);
                $scope.properties.splice(0, $scope.properties.length);
                for (var i = 0; i < tmp.length; i++) {
                    if (tmp[i].id != prevId) {
                        $scope.properties.push(tmp[i]);
                    } else {
                        $scope.properties.push({id: $scope.id, type: $scope.typeSelected.value, desc: $scope.desc});
                    }
                }
            } else {
                $scope.properties.push({id: $scope.id, type: $scope.typeSelected.value, desc: $scope.desc});
            }
            ngDialog.close();
        }
    };
    $scope.goBack = function() {
        ngDialog.close();
    };
});