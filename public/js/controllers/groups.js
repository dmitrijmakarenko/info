accessSettings.controller('groupsCntl', function ($scope, Groups) {
    $scope.groups = [{id: 'group1'}, {id: 'group2'}, {id: 'group3'}];

    $scope.go = function(group) {
        window.location = "#/groups/" + group;
    };

    Groups.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.groups = data.groups||[];
        }
    });
});

accessSettings.controller('groupCntl', function ($scope, $routeParams, ngDialog, Groups) {
    var group = $routeParams.group;

    $scope.createMode = (group != "!new");
    $scope.members = [];
    $scope.users = [];

    Groups.Data.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.users = data.users||[];
            $scope.groups = data.groups||[];
        }
    });

    if (group != "!new") {
        Groups.Get.go({id: group}, function(data) {
            $scope.name = data.name;
            $scope.members = data.members||[];
            $scope.users = data.users||[];
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
        ngDialog.open({
            template: 'addParentDlg',
            controller: 'addParentDlgCntl',
            disableAnimation: true,
            showClose: false,
            scope: $scope
        });
    };

    $scope.removeParent = function(idx) {
        $scope.parents.splice(idx, 1);
    };

    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.name = $scope.name;
        compileSettings.members = $scope.members;
        compileSettings.parents = [];
        for (var i = 0; i < $scope.parents.length; i++) {
            compileSettings.parents.push($scope.parents[i].id);
        }
        Groups.Update.go({id: group, settings: JSON.stringify(compileSettings)}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Отдел сохранен");
                window.location = "#/groups/";
            }
        });
    };

    $scope.deleteGroup = function() {
        Groups.Delete.go({id: group}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Отдел удален");
                window.location = "#/groups/";
            }
        });
    };
});

accessSettings.controller('addParentDlgCntl', function ($scope, Groups, ngDialog) {

    $scope.addParent = function() {
        if ($scope.groupSelected) {
            $scope.parents.push($scope.groupSelected);
            ngDialog.close();
        }
    };

});