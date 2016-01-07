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
            Get: $resource('/get_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id', returnData: '@returnData'} }}),
            Create: $resource('/create_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id', name: '@name', props: '@props'} }}),
            Update: $resource('/update_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id', name: '@name', props: '@props'} }}),
            Remove: $resource('/remove_entity', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

infoSysServices.factory('Config', ['$resource',
    function($resource) {
        return {
            GenerateDB: $resource('/generate_db', {}, { go: { method:'POST', isArray: false }}),
            ValidateDB: $resource('/validate_db', {}, { go: { method:'POST', isArray: false }}),
            ClearDB: $resource('/clear_db', {}, { go: { method:'POST', isArray: false }})
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