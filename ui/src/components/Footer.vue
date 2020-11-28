<template>
  <div :class="isMobile() ? 'container-m' : 'container'">
    <div class="contacts">
      <p>
        Made with
        <img src="@/assets/svg/golang.svg" alt="Go" srcset="" />,
        <img src="@/assets/svg/vue.svg" alt="Vue.js" srcset="" />
        and
        <img src="@/assets/svg/heart.svg" alt="love" srcset="" />
      </p>
    </div>
    <div class="mentions primary">
      <p>
        2020 | by
        <a href="http://github.com/sarmerer" class="secondary">
          Yerzhan Sarmukhanov
        </a>
        &
        <a href="http://github.com/sarmai" class="secondary">
          Ailen Sarmukhanovava
        </a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.container {
  padding-top: 20px;
  margin-top: 30px;
  padding-bottom: 20px;
  color: white;
  /* background-color: rgba(255, 255, 255, 0.03); */
  display: flex;
  flex-direction: column;
  align-items: center;
}

.container-m {
  padding-top: 10px;
  margin-top: 30px;
  padding-bottom: 50px;
  color: white;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.contacts {
  display: flex;
  justify-content: space-around;
  align-items: center;
  vertical-align: middle;
  width: 100%;
}

.mentions p {
  font-size: 0.5em;
}
</style>
