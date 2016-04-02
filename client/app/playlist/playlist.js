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


    }])

    .controller('Playlist', ['$scope', '$http', function ($scope, $http) {


        $scope.videos = [];
        $scope.Date;
        $scope.StartTime = '';
        $scope.LiveSettings = {};
        var afterlive = 0;
        var startv2;
        var entercase = 0

        $scope.format = function (seconds) {
            return new Date(seconds * 1000).toISOString().substr(11, 8);
        };

        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'fakelive/livestreamset'

            }).then(function successCallback(response) {
                console.log(response);
                $scope.LiveSettings = response.data;
                console.log($scope.LiveSettings)
                $scope.LiveSettings.StartTime = string2date(response.data.StartTime)


                $scope.LiveSettings.EndTime = string2date(response.data.EndTime)
                $scope.LiveSettings.StartTime.setSeconds(0)
                $scope.LiveSettings.EndTime.setSeconds(0)
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

            $http({
                method: 'GET',
                url: window.hostname + 'fakelive/getplaylist'
            }).then(function successCallback(response) {
                console.log(response);
                $scope.videos = [];
                $scope.StartTime = response.data.StartTime;
                console.log(response.data.StartTime);
                startv2 = new Date($scope.StartTime);

                var once = 0;

                angular.forEach(response.data.Videos, function (value, key) {

                    var res = calcVideoTime(value.Duration);
                    if ($scope.LiveSettings.Activated) {
                        //if video ends after live start
                        if (res[2].getTime() > $scope.LiveSettings.StartTime.getTime()) {

                            //enter once
                            if (afterlive < 1) {
                                afterlive++;
                                var cuttime = res[2].getTime() - $scope.LiveSettings.StartTime.getTime()
                                var stoptime = value.Duration - cuttime/1000
                                console.log(value.Duration,$scope.LiveSettings.StartTime)
                                value.Duration = stoptime
                                res = calcVideoTime(stoptime);
                                value.PlayingInterval = res[0];
                                value.Playing = dateCheck(res[1], res[2], new Date());
                                this.push(value);
                                sumseconds = (sumseconds + Number(stoptime));

                                var livetime = ($scope.LiveSettings.EndTime.getTime() - $scope.LiveSettings.StartTime.getTime()) / 1000
                                res = calcVideoTime(livetime);
                                value = angular.copy(value);
                                value.PlayingInterval = res[0];
                                value.Thumbnail='http://www.azorestv.com/uploads/images/clip_410_1428943639_poster.jpg';
                                value.Title="Live Streaming"
                                value.Id=25;
                                value.Duration=livetime
                                value.Playing = dateCheck(res[1], res[2], new Date());

                                sumseconds = (sumseconds + Number(livetime));
                                this.push(value);

                                once=1


                            }else{
                                once++
                            }

                        }

                    }

                    if (once != 1) {
                        once=2

                        sumseconds = (sumseconds + Number(value.Duration));
                        value.PlayingInterval = res[0];
                        value.Playing = dateCheck(res[1], res[2], new Date());
                        this.push(value);

                    }


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

        var startHour = 18;

        var sumseconds = 0;

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

        var calcVideoTime = function (seconds) {
            var intseconds = seconds * 1;

            var datestart = new Date(startv2.getTime() + (sumseconds * 1000));

            var dateend = new Date(datestart.getTime() + (seconds) * 1000);

            var res = [];
            res[0] = datetoHHMMSS(datestart) + ' - ' + datetoHHMMSS(dateend);
            res[1] = datestart;
            res[2] = dateend;
            return res
        };


    }]).controller('Settings', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {


    $scope.FakeliveSettings={
        LiveStreamSettings:  {StartLiveTime: new Date(), EndLiveTime: new Date()},
        StartTime:'',
        RTimes:[]
    };

    $scope.removeRTime=function(index){
        $scope.FakeliveSettings.RTimes.splice(index, 1);
    };



    $scope.addRepeatTimes=function(){

        $scope.FakeliveSettings.RTimes.push(new Date());

    };


    var update = function () {
        $http({
            method: 'GET',
            url: window.hostname + 'fakelive/settings'
        }).then(function successCallback(response) {
            $scope.FakeliveSettings=response.data;

            $scope.FakeliveSettings.StartTime = string2date($scope.FakeliveSettings.StartTime);
            $scope.FakeliveSettings.LiveStreamSettings.StartLiveTime = string2date($scope.FakeliveSettings.LiveStreamSettings.StartLiveTime)
            $scope.FakeliveSettings.LiveStreamSettings.EndLiveTime = string2date($scope.FakeliveSettings.LiveStreamSettings.EndLiveTime)
            if ($scope.FakeliveSettings.RTimes ==null) {


                $scope.FakeliveSettings.RTimes=[];

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

        $scope.FakeliveSettings.StartTime=date2string($scope.FakeliveSettings.StartTime);
        $scope.FakeliveSettings.LiveStreamSettings.StartLiveTime=date2string($scope.FakeliveSettings.LiveStreamSettings.StartLiveTime);
        $scope.FakeliveSettings.LiveStreamSettings.EndLiveTime=date2string($scope.FakeliveSettings.LiveStreamSettings.EndLiveTime);

        $http({
            method: 'POST',
            url: window.hostname + 'fakelive/settings',
            data:$scope.FakeliveSettings,
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





}]);