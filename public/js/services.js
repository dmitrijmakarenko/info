var infoSysServices = angular.module('infoSysServices', ['ngResource']);

infoSysServices.factory('informerControl', function() {
    return {
        status: null,
        message: null,
        success: function(msg) {
            this.status = 'success';
            this.message = msg;
        },
        error: function(msg) {
            this.status = 'error';
            this.message = msg;
        },
        clear: function() {
            this.status = null;
            this.message = null;
        }
    }
});

infoSysServices.factory('Entity', ['$resource',
    function($resource) {
        return {
            List: $resource('/entities', {}, { go: { method:'POST', isArray: true }}),
            Get: $resource('/get_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
            //Create: $resource('/create_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id', name: '@name', props: '@props'} }}),
            //Update: $resource('/update_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id', name: '@name', props: '@props'} }}),
            //Remove: $resource('/remove_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

infoSysServices.factory('DataBase', ['$resource',
    function($resource) {
        return {
            GetTables: $resource('/tables', {}, { go: { method:'POST', isArray: true }}),
            GetViews: $resource('/views', {}, { go: { method:'POST', isArray: true }}),
            Protect: $resource('/protect', {}, { go: { method:'POST', isArray: false, params : {table: '@table'} }})
        };
    }]);

infoSysServices.factory('Accounts', ['$resource',
    function($resource) {
        return {
            List: $resource('/users/list', {}, { go: { method:'POST', isArray: false }}),
            Get: $resource('/users/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }}),
            Update: $resource('/users/update', {}, { go: { method:'POST', isArray: false, params : {id: '@id', data: '@settings'} }}),
            Delete: $resource('/users/delete', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

infoSysServices.factory('Groups', ['$resource',
    function($resource) {
        return {
            List: $resource('/groups/list', {}, { go: { method:'POST', isArray: false }}),
            Get: $resource('/groups/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }}),
            Update: $resource('/groups/update', {}, { go: { method:'POST', isArray: false, params : {id: '@id', name: '@name'} }}),
            Delete: $resource('/groups/delete', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

infoSysServices.factory('Rules', ['$resource',
    function($resource) {
        return {
            List: $resource('/rules', {}, { go: { method:'POST', isArray: false }})
            //Get: $resource('/account', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

infoSysServices.factory('TestFunc', ['$resource',
    function($resource) {
        return {
            Select: $resource('/testdata', {}, { go: { method:'POST', isArray: false }})
        };
    }]);