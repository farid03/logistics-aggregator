let myMap;

// Подробнее о действиях в https://yandex.ru/dev/maps/jsbox/2.1/icon_customImage
ymaps.ready(init);

function init () {
    coords = document.getElementsByClassName("car")
    myMap = new ymaps.Map('map', {
        center: [55.76, 37.64], // Москва
        zoom: 10
    }, {
        searchControlProvider: 'yandex#search'
    });

    for (let i = 0; i < coords.length; i++) {
        if (Number(coords[i].dataset.latitude) === 0 &&
            Number(coords[i].dataset.longitude) === 0) {
            continue
        }
        addCarOnMap(coords[i].dataset)
    }
}

function addCarOnMap(dataset) {
    let carPlacemark = new ymaps.Placemark([Number(dataset.latitude), Number(dataset.longitude)], {
        hintContent: 'Регистрационный номер: ' + dataset.licenseplate,
        balloonContent: 'Широта: ' + dataset.latitude + "\n" + "Долгота: " + dataset.longitude
    }, {
        iconLayout: 'default#image',
        iconImageHref: 'static/img/gazon.png',
        iconImageSize: [100, 70],
        iconImageOffset: [-50, -35],
    });
    myMap.geoObjects.add(carPlacemark)
}