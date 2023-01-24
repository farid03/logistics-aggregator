ymaps.ready({
    successCallback: () => {
        coords = document.getElementsByClassName("geocode")

        for (let i = 0; i < coords.length; i++) {
            coordArgs = coords[i].innerHTML.trim().split(" ")
            latitude = Number(coordArgs[1])
            longitude = Number(coordArgs[2].substring(0, coordArgs[2].length - 1))
            if (latitude === 0 && longitude === 0) {
                coords[i].innerHTML = "Not defined"
                continue
            }
            ymaps.geocode([latitude, longitude]).then(
                function (res) {
                    address = res.geoObjects.get(0).properties.getAll().text
                    coords[i].innerHTML = address
                },
                function (err) {
                    console.log(err)
                }
            )
        }
    }
})