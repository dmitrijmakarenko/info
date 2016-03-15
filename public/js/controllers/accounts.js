accessSettings.controller('accountsCntl', function ($scope, Accounts) {

    $scope.go = function(account) {
        window.location = "#/accounts/" + account;
    };

    Accounts.List.go(function(data) {
        $scope.accounts = [];
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            for (var i = 0; i < data.accounts.length; i ++) {
                var item = {};
                item.id = data.accounts[i].id;
                item.name = data.accounts[i].name||data.accounts[i].id;
                $scope.accounts.push(item);
            }
        }
    });
});

accessSettings.controller('accountCntl', function ($scope, $routeParams, Accounts) {
    var account = $routeParams.account;

    $scope.createMode = (account != "!new");

    if (account != "!new") {
        Accounts.Get.go({id: account}, function(data) {
            console.log("get user", data);
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.id = account;
                $scope.name = data.name;
                $scope.position = data.position;
            }
        });
    }

    $scope.saveSettings = function() {
        var compileSettings = {};
        compileSettings.id = $scope.id;
        compileSettings.name = $scope.name;
        compileSettings.position = $scope.position;
        Accounts.Update.go({id: account, settings: JSON.stringify(compileSettings)}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Пользователь сохранен");
                window.location = "#/accounts/";
            }
        });
    };

    $scope.deleteAccount = function() {
        Accounts.Delete.go({id: account}, function(data) {
            if (data.error) {
                $scope.showErrorMsg(data.error);
            } else {
                $scope.showSuccessMsg("Пользователь удален");
                window.location = "#/accounts/";
            }
        });
    };
});