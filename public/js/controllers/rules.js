infoSys.controller('rulesCntl', function ($scope, Rules) {
    $scope.rules = [{id: 'dasdasdasd'}, {id: 'as21xczxsdsa'}, {id: '12312zxcz'}];

    Rules.List.go(function(data) {
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.rules = data.rules||[];
        }
    });

});
