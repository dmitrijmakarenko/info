accessSettings.controller('mainCntl', function ($scope, informerControl, spinnerControl) {
    $scope.showSuccessMsg = function(msg) { informerControl.success(msg); };
    $scope.showErrorMsg = function(msg) { informerControl.error(msg); };

    $scope.showSpinner = function() { spinnerControl.showSpinner(); };
    $scope.hideSpinner = function() { spinnerControl.hideSpinner(); }
});