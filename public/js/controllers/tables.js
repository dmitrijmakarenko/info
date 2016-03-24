accessSettings.controller('tablesCntl', function ($scope, Entity) {
    Entity.List.go(function(data) {
        $scope.entities = data;
    });

    $scope.getEntity = function(entity) {
        window.location = "#/tables/" + entity;
    };

});

accessSettings.controller('tableCntl', function ($scope, $routeParams, ngDialog, Entity) {
    //$scope.loading = true;
    var table = $routeParams.tableId,
        uuids = [];

    Entity.Get.go({id: table}, function(data) {
        //console.log("table data", data);
        if (data.error) {
            $scope.showErrorMsg(data.error);
        } else {
            $scope.entity = data;
            if (!$scope.entity.rows) $scope.entity.rows = [];
            uuids = data.uuid||[];
        }
    });

    $scope.protectRec = function(idx) {
        if (uuids.length > 0) {
            $scope.table = table;
            $scope.uuid = uuids[idx];
            ngDialog.open({
                template: 'protectRecDlg',
                controller: 'protectRecDlgCntl',
                disableAnimation: true,
                showClose: false,
                scope: $scope
            });
        }
    };

});

accessSettings.controller('protectRecDlgCntl', function ($scope, ngDialog, Rules) {
    $scope.rules = [];
    Rules.List.go(function(data) {
        if (!data.error) {
            $scope.rules = data.rules||[];
        }
    });
    $scope.protect = function() {
        if ($scope.uuid && $scope.ruleSelected && $scope.ruleSelected.id) {
            var settings = {};
            settings.uuidRecord = $scope.uuid;
            settings.uuidRule = $scope.ruleSelected.id;
            settings.table = $scope.table;
            Rules.ProtectRecord.go({settings: settings}, function(data) {
                if (!data.error) {
                    ngDialog.close();
                }
            });
        }
    };
});
