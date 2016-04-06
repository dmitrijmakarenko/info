accessSettings.controller('mainCntl', function ($scope, informerControl) {
    $scope.showSuccessMsg = function(msg) { informerControl.success(msg); };
    $scope.showErrorMsg = function(msg) { informerControl.error(msg); }
});

accessSettings.controller('mainCntl', function ($scope, spinnerControl) {
    $scope.showSpinner = function() { spinnerControl.showSpinner(); };
    $scope.hideSpinner = function() { spinnerControl.hideSpinner(); }
});