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

infoSys.controller('accountCntl', function ($scope, $routeParams, Accounts) {
    var account = $routeParams.account;

    $scope.createMode = (account != "!new");

    if (account != "!new") {
        Accounts.Get.go({id: account}, function(data) {
            console.log(data);
            $scope.id = account;
            $scope.name = data.name;
            $scope.position = data.position;
        });
    }

    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.id = $scope.id;
        compileSettings.name = $scope.name;
        compileSettings.position = $scope.position;
        Accounts.Update.go({create: (account == "!new"), settings: JSON.stringify(compileSettings)}, function(data) {
            console.log(data);
        });
    }
});