'use strict';
function string2date(strdate) {

    if (!strdate)
        return new Date();

    if (strdate.indexOf(":") > -1) {
        var date = strdate.split(":");
        var hour, minute;
        hour = date[0];
        minute = date[1];
        var d = new Date();
        d.setHours(hour);
        d.setMinutes(minute);
        return d
    }


}
angular.module('myApp.Playlist', ['ngRoute'])
    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/playlist/list', {
            templateUrl: window.hostnametpl + 'playlist/list.html',
            controller: 'Playlist'
        });


        $routeProvider.when('/playlist/settings', {
            templateUrl: window.hostnametpl + 'playlist/settings.html',
            controller: 'Settings'
        });

        $routeProvider.when('/playlist/manage', {
            templateUrl: window.hostnametpl + 'playlist/manage.html',
            controller: 'Manage'
        });


    }])

    .controller('Playlist', ['$scope', '$http', function ($scope, $http) {


        $scope.videos = [];
        $scope.Date;


        $scope.format = function (seconds) {
            return new Date(seconds * 1000).toISOString().substr(11, 8);
        };

        var update = function () {


            $http({
                method: 'GET',
                url: window.hostname + 'fakelive/getplaylist'
            }).then(function successCallback(response) {

                $scope.videos = [];


                angular.forEach(response.data, function (value, key) {


                    value.PlayingInterval = datetoHHMMSS(new Date(value.Scheduled)) + ' - ' + datetoHHMMSS(new Date(value.EndTime));
                    value.Duration = (new Date(value.EndTime).getTime() - new Date(value.Scheduled).getTime()) / 1000
                    value.Playing = dateCheck(new Date(value.Scheduled), new Date(value.EndTime), new Date());
                    this.push(value);


                }, $scope.videos);

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update();





        function dateCheck(from, to, check) {
            if ((check.getTime() <= to.getTime() && check.getTime() >= from.getTime())) {
                return true;
            }
            return false;
        }


        function datetoHHMMSS(date) {
            return (date.getHours() < 10 ? "0" + date.getHours() : date.getHours()) +
                ":" + (date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes()) +
                ":" + (date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds());
        }



    }]).controller('Settings', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {


    $scope.FakeliveSettings = {
        LiveStreamSettings: {StartLiveTime: new Date(), EndLiveTime: new Date()},
        StartTime: '',
        RTimes: []
    };

    $scope.removeRTime = function (index) {
        $scope.FakeliveSettings.RTimes.splice(index, 1);
    };


    $scope.addRepeatTimes = function () {

        $scope.FakeliveSettings.RTimes.push(new Date());

    };


    var update = function () {
        $http({
            method: 'GET',
            url: window.hostname + 'fakelive/settings'
        }).then(function successCallback(response) {
            $scope.FakeliveSettings = response.data;

            $scope.FakeliveSettings.StartTime = string2date($scope.FakeliveSettings.StartTime);
            $scope.FakeliveSettings.LiveStreamSettings.StartLiveTime = string2date($scope.FakeliveSettings.LiveStreamSettings.StartLiveTime)
            $scope.FakeliveSettings.LiveStreamSettings.EndLiveTime = string2date($scope.FakeliveSettings.LiveStreamSettings.EndLiveTime)
            if ($scope.FakeliveSettings.RTimes == null) {
                $scope.FakeliveSettings.RTimes = [];

            }


            // $scope.StartTime=response.data.StartTime

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });


    };
    update();

    function pad(n, width, z) {
        z = z || '0';
        n = n + '';
        return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
    }

    function date2string(date) {
        return pad(date.getHours(), 2) + ":" + pad(date.getMinutes(), 2)
    }

    $scope.reload = function () {

        $http({
            method: 'POST',
            url: window.hostname + 'fakelive/reload',
        }).then(function successCallback(response) {
            toastr.success('Success!', 'servidor reiniciado');
            console.log(response);
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

    };


    $scope.updateSettings = function () {

        $scope.FakeliveSettings.StartTime = date2string($scope.FakeliveSettings.StartTime);
        $scope.FakeliveSettings.LiveStreamSettings.StartLiveTime = date2string($scope.FakeliveSettings.LiveStreamSettings.StartLiveTime);
        $scope.FakeliveSettings.LiveStreamSettings.EndLiveTime = date2string($scope.FakeliveSettings.LiveStreamSettings.EndLiveTime);

        $http({
            method: 'POST',
            url: window.hostname + 'fakelive/settings',
            data: $scope.FakeliveSettings,
        }).then(function successCallback(response) {
            console.log(response);
            toastr.success('Success!', 'Definicoes guardadas');
            update();

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });


    };


}]).controller('Manage', ['$scope', '$http','toastr', function ($scope, $http,toastr) {


    $scope.videos = [];





    var update = function () {


        $http({
            method: 'GET',
            url: window.hostname + 'fakelive/trim/new'
        }).then(function successCallback(response) {

            $scope.videos = response.data;

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    };
    update();

    $scope.updateVideo=function(video){
        $http({
            method: 'POST',
            url: window.hostname + 'fakelive/trim/save',
            data:video
        }).then(function successCallback(response) {

            toastr.success('Success!', 'Video actualizado');

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Success!', response.data);
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

    };

    $scope.search=function (keywords){
        if (keywords==""){
            update()
            return

        }
        $http({
            method: 'GET',
            url: window.hostname + 'fakelive/trim/search/'+keywords
        }).then(function successCallback(response) {

            $scope.videos = response.data;

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });


    }




}]);