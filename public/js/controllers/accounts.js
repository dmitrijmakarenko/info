infoSys.controller('accountsCntl', function ($scope) {
    $scope.accounts = [{id: 'user1'}, {id: 'user2'}, {id: 'user3'}];

    $scope.go = function(account) {
        window.location = "#/accounts/" + account;
    }
});

infoSys.controller('accountCntl', function ($scope, $routeParams) {
    var account = $routeParams.account;

    if (account != "!new") {
        $scope.createMode = true;
    } else {
        $scope.createMode = false;
    }
});