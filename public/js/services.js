var accessSettingsServices = angular.module('accessSettingsServices', ['ngResource']);

accessSettingsServices.factory('informerControl', function() {
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

accessSettingsServices.factory('spinnerControl', function() {
    return {
        spinnerShow: false,
        showSpinner: function() {
            this.spinnerShow = true;
        },
        hideSpinner: function() {
            this.spinnerShow = false;
        }
    }
});

accessSettingsServices.factory('Tables', ['$resource',
    function($resource) {
        return {
            List: $resource('/entities/list', {}, { go: { method:'POST', isArray: true }}),
            Get: $resource('/tables/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

accessSettingsServices.factory('DataBase', ['$resource',
    function($resource) {
        return {
            GetTables: $resource('/tables/list', {}, { go: { method:'POST', isArray: true }}),
            GetViews: $resource('/views', {}, { go: { method:'POST', isArray: true }}),
            Protect: $resource('/protect', {}, { go: { method:'POST', isArray: false, params : {table: '@table'} }})
        };
    }]);

accessSettingsServices.factory('Accounts', ['$resource',
    function($resource) {
        return {
            List: $resource('/users/list', {}, { go: { method:'POST', isArray: false }}),
            Get: $resource('/users/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }}),
            Update: $resource('/users/update', {}, { go: { method:'POST', isArray: false, params : {id: '@id', data: '@settings'} }}),
            Delete: $resource('/users/delete', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

accessSettingsServices.factory('Groups', ['$resource',
    function($resource) {
        return {
            List: $resource('/groups/list', {}, { go: { method:'POST', isArray: false }}),
            Get: $resource('/groups/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }}),
            Data: $resource('/groups/data', {}, { go: { method:'POST', isArray: false }}),
            Update: $resource('/groups/update', {}, { go: { method:'POST', isArray: false, params : {id: '@id', data: '@settings'} }}),
            Delete: $resource('/groups/delete', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

accessSettingsServices.factory('Rules', ['$resource',
    function($resource) {
        return {
            List: $resource('/rules/list', {}, { go: { method:'POST', isArray: false }}),
            Get: $resource('/rules/get', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }}),
            Data: $resource('/rules/data', {}, { go: { method:'POST', isArray: false }}),
            Update: $resource('/rules/update', {}, { go: { method:'POST', isArray: false, params : {id: '@id', data: '@settings'} }}),
            ProtectRecord: $resource('/rules/protectrec', {}, { go: { method:'POST', isArray: false, params : {data: '@settings'} }}),
            Delete: $resource('/rules/delete', {}, { go: { method:'POST', isArray: false, params : {id: '@id'} }})
        };
    }]);

accessSettingsServices.factory('Test', ['$resource',
    function($resource) {
        return {
            Auth: $resource('/auth', {}, { go: { method:'POST', isArray: false, params : {user: '@user', pass: '@pass'} }}),
            Install: $resource('/test/install', {}, { go: { method:'POST', isArray: false }}),
            Reset: $resource('/test/reset', {}, { go: { method:'POST', isArray: false }}),
            Init: $resource('/test/init', {}, { go: { method:'POST', isArray: false }}),
            Work: $resource('/test/work', {}, { go: { method:'POST', isArray: false }}),
            Compile: $resource('/test/compile', {}, { go: { method:'POST', isArray: false }}),
            CopyToFile: $resource('/test/copytofile', {}, { go: { method:'POST', isArray: false }}),
            CopyFromFile: $resource('/test/copyfromfile', {}, { go: { method:'POST', isArray: false }}),
            SelectDataNormal: $resource('/test/selectdatanormal', {}, { go: { method:'POST', isArray: false }}),
            SelectDataSecure: $resource('/test/selectdatasecure', {}, { go: { method:'POST', isArray: false }}),
            UpdateDataNormal: $resource('/test/updatedatanormal', {}, { go: { method:'POST', isArray: false }}),
            UpdateDataSecure: $resource('/test/updatedatasecure', {}, { go: { method:'POST', isArray: false }}),
            DeleteDataNormal: $resource('/test/deletedatanormal', {}, { go: { method:'POST', isArray: false }}),
            DeleteDataSecure: $resource('/test/deletedatasecure', {}, { go: { method:'POST', isArray: false }})
        };
    }]);

accessSettingsServices.factory('VCS', ['$resource',
    function($resource) {
        return {
            Tables: $resource('/vcs/tables', {}, { go: { method:'POST', isArray: false }}),
            AddToVcs: $resource('/vcs/add', {}, { go: { method:'POST', isArray: false, params : {table: '@table'} }}),
            RemoveFromVcs: $resource('/vcs/delete', {}, { go: { method:'POST', isArray: false, params : {table: '@table'} }})
        };
    }]);

accessSettingsServices.factory('Data', ['$resource',
    function($resource) {
        return {
            Get: $resource('/data/get', {}, { go: { method:'POST', isArray: false, params : {params: '@params'} }}),
            Add: $resource('/data/add', {}, { go: { method:'POST', isArray: false, params : {params: '@params'} }})
        };
    }]);