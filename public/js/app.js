var accessSettings = angular.module('accessSettings', ['ngRoute', 'ngAnimate', 'accessSettingsServices', 'accessSettingsDirectives', 'ngDialog']);

accessSettings.run(function($rootScope) {
    angular.element(document).on("click", function(e) {
        $rootScope.$broadcast("documentClicked", angular.element(e.target));
    });
});

accessSettings.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
            when('/rules', {
                templateUrl: 'public/templates/rules.html',
                controller: 'rulesCntl'
            }).
            when('/rules/:rule', {
                templateUrl: 'public/templates/rule.html',
                controller: 'ruleCntl'
            }).
            when('/tables', {
                templateUrl: 'public/templates/tables.html',
                controller: 'tablesCntl'
            }).
            when('/table/:tableId', {
                templateUrl: 'public/templates/table.html',
                controller: 'tableCntl'
            }).
            when('/accounts', {
                templateUrl: 'public/templates/accounts.html',
                controller: 'accountsCntl'
            }).
            when('/accounts/:account', {
                templateUrl: 'public/templates/account.html',
                controller: 'accountCntl'
            }).
            when('/groups', {
                templateUrl: 'public/templates/groups.html',
                controller: 'groupsCntl'
            }).
             when('/groups/:group', {
                templateUrl: 'public/templates/group.html',
                controller: 'groupCntl'
            }).
            when('/test', {
                templateUrl: 'public/templates/test.html',
                controller: 'testCntl'
            }).
            when('/db', {
                templateUrl: 'public/templates/db.html',
                controller: 'dbCntl'
            }).
            otherwise({
                redirectTo: '/db'
            });
    }]);