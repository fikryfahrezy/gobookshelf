/**
 * @typedef {{alpha2Code: string, name: string}} Country
 */
const form = document.getElementById('form');
const regionInput = document.getElementById('region-input');
const regionSelect = document.getElementById('region-select');
const streetInput = document.getElementById('street-input');
const streetSelect = document.getElementById('street-select');

let selectedRegion = regionInput.value;

/**
 *
 * @param {boolean} isSelected
 */
const addressVisibleToggle = function addressVisible(isSelected) {
  if (isSelected) {
    streetInput.disabled = false;
    streetSelect.disabled = false;
  } else {
    streetInput.disabled = true;
    streetSelect.disabled = true;
  }
};

/**
 *
 * @param {Country} param0
 */
const appendCountry = function appendCountry({ alpha2Code, name }) {
  const option = document.createElement('option');

  option.classList.add('country');
  option.value = alpha2Code;
  option.textContent = name;

  option.addEventListener('click', (e) => {
    selectedRegion = e.target.value;
    regionInput.value = e.target.textContent;

    addressVisibleToggle(selectedRegion !== '');
  });

  regionSelect.appendChild(option);
};

/**
 *
 * @param {Country[]} countries
 */
const appendCountries = function setCountries(countries) {
  while (regionSelect.lastElementChild) {
    regionSelect.removeChild(regionSelect.lastElementChild);
  }

  countries.forEach((country) => {
    appendCountry(country);
  });
};

const setDefaultCountry = function setDefalutCountry() {
  const option = document.querySelector('.address-input  .country');

  if (option) {
    option.click();
  }
};

/**
 *
 * @param {string} street
 */
const appendStreet = function appendCountry(street) {
  const option = document.createElement('option');

  option.classList.add('street');
  option.value = street;
  option.textContent = street;

  option.addEventListener('click', (e) => {
    streetInput.value = e.target.textContent;
  });

  streetSelect.appendChild(option);
};

/**
 *
 * @param {string[]} streets
 */
const appendStreets = function setCountries(streets) {
  while (streetSelect.lastElementChild) {
    streetSelect.removeChild(streetSelect.lastElementChild);
  }

  streets.forEach((street) => {
    appendStreet(street);
  });
};

/**
 *
 * @param {string} name
 * @returns {Promise<Country[]>}
 */
const fetchCountries = async function fetchCountries(name = '') {
  return fetch(`/countries?name=${name}`)
    .then((res) => res.json())
    .then((res) => {
      if (!Array.isArray(res)) throw new Error('not found');

      const data = res.map(({ alpha2Code, name }) => ({ alpha2Code, name }));

      return data;
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 *
 * @param {string} street
 * @returns {Promise<string[]>}
 */
const fetchStreet = async function fetchStreet(street = 'jakarta') {
  return fetch(`/street?region=${selectedRegion}&street=${street}`, {
    method: 'GET',
  })
    .then((res) => res.json())
    .then(({ standard }) => {
      if (!standard) throw new Error('not-found');

      const data = Object.keys(standard.street);

      return data;
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 *
 * @param {string} name
 */
const getCountries = async function getCountries(name) {
  return fetchCountries(name)
    .then((res) => {
      if (!res) throw new Error('not-found');

      appendCountries(res);
    })
    .then(() => {
      setDefaultCountry();
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 *
 * @param {string} street
 */
const getStreet = async function getStreet(street) {
  return fetchStreet(street)
    .then((res) => {
      if (!res) throw new Error('not-found');

      appendStreets(res);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 *
 * @param {(e: Event) => void} fn
 * @param {number} timer
 * @returns {(e: Event) => void}
 */
const eventDebouncer = function eventDebouncer(fn, timer = 300) {
  let timeout = setTimeout(() => {}, 0);
  return function (e) {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      fn(e);
    }, timer);
  };
};

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const form = new FormData(e.target);
  const data = Object.fromEntries(form.entries());

  fetch('/registration', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (res.status >= 200 && res.status < 300) {
        window.location = '/';
      }
    })
    .catch((err) => {
      console.log(err);
    });
});

regionInput.addEventListener(
  'keydown',
  eventDebouncer((e) => {
    getCountries(e.target.value);
  }, 300)
);

regionInput.addEventListener(
  'focusin',
  eventDebouncer((e) => {
    getCountries(e.target.value);
  }, 300)
);

streetInput.addEventListener(
  'keydown',
  eventDebouncer((e) => {
    getStreet(e.target.value);
  }, 300)
);

streetInput.addEventListener(
  'focusin',
  eventDebouncer((e) => {
    getStreet(e.target.value);
  }, 300)
);

const init = async function init() {
  addressVisibleToggle(selectedRegion !== '');
  await getCountries(selectedRegion);
};

init();
