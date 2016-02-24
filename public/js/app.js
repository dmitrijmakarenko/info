var infoSys = angular.module('infoSys', ['ngRoute', 'ngAnimate', 'infoSysServices', 'infoSysDirectives', 'ngDialog']);

infoSys.run(function($rootScope) {
    angular.element(document).on("click", function(e) {
        $rootScope.$broadcast("documentClicked", angular.element(e.target));
    });
});

infoSys.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
            /*when('/objects', {
                templateUrl: 'public/templates/objects.html',
                controller: 'objectsCntl'
            }).
            when('/objects/:objectId', {
                templateUrl: 'public/templates/object.html',
                controller: 'objectCntl'
            }).
            when('/objects-edit', {
                templateUrl: 'public/templates/objects_edit.html',
                controller: 'entitiesCntl'
            }).
            when('/objects-edit/:objectId', {
                templateUrl: 'public/templates/object_edit.html',
                controller: 'entityCntl'
            }).
            when('/configs', {
                templateUrl: 'public/templates/configs.html',
                controller: 'configsCntl'
            }).*/
            when('/rules', {
                templateUrl: 'public/templates/rules.html',
                controller: 'rulesCntl'
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
            when('/db', {
                templateUrl: 'public/templates/db.html',
                controller: 'dbCntl'
            }).
            otherwise({
                redirectTo: '/db'
            });
    }]);