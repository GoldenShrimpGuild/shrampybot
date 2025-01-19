<template>
  <RouterView />
</template>

<style lang="scss">
* {
  font-family: 'Rajdhani';
}

#app {
  font-family: 'Rajdhani';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #ffffff;
}

body {
  margin: 0;
}

div.va-card {
  &__title {
    align-items: top;
    margin-left: 1em;
    font-family: 'Revalia';
    font-size: 10pt;
  }
}

h1 {
  font-size: 24pt;
}
h2 {
  font-size: 22pt;
}
h3 {
  font-size: 20pt;
}
h4 {
  font-size: 16pt;
}
h5 {
  font-size: 14pt;
}

h1,
h2,
h3,
h4,
h5 {
  font-family: 'Revalia';
  text-transform: uppercase;
}

.page-title {
  font-size: 24pt;
}

.card-title {
  font-size: 16pt;
}

a,
div {
  span.logo {
    font-family: 'Revalia';
    font-size: 20pt;
    color: #ffbb22;
    font-weight: bold;
    .dot_text {
      font-family: 'Revalia';
      font-size: 20pt;
      color: #e42222;
    }
    .bot_text {
      font-family: 'Revalia';
      font-size: 20pt;
      color: #ffffff;
    }
    .caridea_text {
      font-family: 'Revalia';
      font-size: 20pt;
    }
  }
}
a {
  color: #ffbb22;
}
</style>
