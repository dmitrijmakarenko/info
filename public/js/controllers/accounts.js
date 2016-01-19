infoSys.controller('accountsCntl', function ($scope, Accounts) {
    //$scope.accounts = [{id: 'user1'}, {id: 'user2'}, {id: 'user3'}];

    $scope.go = function(account) {
        window.location = "#/accounts/" + account;
    };

    Accounts.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.accounts = data.accounts||[];
        }
    });
});

infoSys.controller('accountCntl', function ($scope, $routeParams) {
    var account = $routeParams.account;

    $scope.createMode = (account != "!new");

    $scope.saveSettings = function() {
        if (account == "!new") {

        } else {

        }
    }
});