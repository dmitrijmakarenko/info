infoSys.controller('groupsCntl', function ($scope, Groups) {
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

infoSys.controller('groupCntl', function ($scope, $routeParams) {
    var group = $routeParams.group;

    $scope.createMode = (group != "!new");

    $scope.saveSettings = function() {
        if (group == "!new") {

        } else {

        }
    }
});