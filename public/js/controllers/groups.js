accessSettings.controller('groupsCntl', function ($scope, Groups) {
    $scope.groups = [{id: 'group1'}, {id: 'group2'}, {id: 'group3'}];

    $scope.go = function(group) {
        window.location = "#/groups/" + group;
    };

    $scope.showSpinner();
    Groups.List.go(function(data) {
        $scope.hideSpinner();
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.groups = data.groups||[];
        }
    });
});

accessSettings.controller('groupCntl', function ($scope, $routeParams, ngDialog, Groups) {
    var group = $routeParams.group,
        groupNameById = {},
        oldSettings = {};

    $scope.createMode = (group != "!new");
    $scope.members = [];
    $scope.users = [];

    $scope.showSpinner();
    Groups.Data.go(function(data) {
        $scope.hideSpinner();
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.users = data.users||[];
            $scope.groups = data.groups||[];
            for (var i = 0; i < $scope.groups.length; i++) {
                if ($scope.groups[i].id) groupNameById[$scope.groups[i].id] = $scope.groups[i].name||$scope.groups[i].id;
            }
        }
    });

    if (group != "!new") {
        Groups.Get.go({id: group}, function(data) {
            $scope.name = data.name;
            $scope.members = data.members||[];
            $scope.users = data.users||[];
            $scope.parents = data.parents||[];
            if (data.parents instanceof Array && data.parents.length > 0) {
                oldSettings.parents = data.parents.slice();
            } else {
                oldSettings.parents = [];
            }
            if ($scope.groups instanceof Array && $scope.groups.length > 0) {

            }
        });
    }

    $scope.addMember = function(id) {
        var tmp = [], item;
        for (var i = 0; i < $scope.users.length; i++) {
            if ($scope.users[i].id != id) {
                tmp.push($scope.users[i]);
            } else {
                item = $scope.users[i];
            }
        }
        if (item) {
            $scope.users = tmp;
            $scope.members.push(item);
        }
    };

    $scope.removeMember = function(id) {
        var tmp = [], item;
        for (var i = 0; i < $scope.members.length; i++) {
            if ($scope.members[i].id != id) {
                tmp.push($scope.members[i]);
            } else {
                item = $scope.members[i];
            }
        }
        if (item) {
            $scope.members = tmp;
            $scope.users.push(item);
        }
    };

    $scope.parents = [];

    $scope.showAddParent = function() {
        $scope.groupsAvailable = [];
        for (var i = 0; i < $scope.groups.length; i++) {
            if ($scope.groups[i].id != group) {
                $scope.groupsAvailable.push($scope.groups[i]);
            }
        }
        ngDialog.open({
            template: 'addParentDlg',
            controller: 'addParentDlgCntl',
            disableAnimation: true,
            showClose: false,
            scope: $scope
        });
    };

    $scope.getName = function(id) {
        return groupNameById[id]||"";
    };

    $scope.removeParent = function(idx) {
        $scope.parents.splice(idx, 1);
    };

    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.name = $scope.name;
        compileSettings.members = $scope.members;
        compileSettings.parentsAdd = [];
        for (var i = 0; i < $scope.parents.length; i++) {
            if ($scope.parents[i].level == 1) {
                var isNew = true;
                for (var j = 0; j < oldSettings.parents.length; j++) {
                    if ($scope.parents[i].id == oldSettings.parents[j].id) {
                        isNew = false;
                        oldSettings.parents.splice(j, 1);
                        break;
                    }
                }
                if (isNew) {
                    var item = {};
                    item.id = $scope.parents[i].id;
                    item.level = 1;
                    compileSettings.parentsAdd.push(item);
                }
            }
        }
        compileSettings.parentsRemove = oldSettings.parents||[];
        Groups.Update.go({id: group, settings: JSON.stringify(compileSettings)}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Группа сохранена");
                window.location = "#/groups/";
            }
        });
    };

    $scope.deleteGroup = function() {
        Groups.Delete.go({id: group}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Группа удалена");
                window.location = "#/groups/";
            }
        });
    };
});

accessSettings.controller('addParentDlgCntl', function ($scope, Groups, ngDialog) {

    $scope.addParent = function() {
        if ($scope.groupSelected && $scope.groupSelected.id) {
            var item = {};
            item.id = $scope.groupSelected.id;
            item.level = 1;
            $scope.parents.push(item);
            ngDialog.close();
        }
    };

});