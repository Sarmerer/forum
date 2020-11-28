<template>
  <div :class="isMobile() ? 'container-m' : 'container'">
    <div class="contacts">
      <p class="made-with">
        Made with
        <a target="_blank" rel="noopener noreferrer" href="https://golang.org">
          <img src="@/assets/svg/golang.svg" alt="Go"/></a
        >,
        <a target="_blank" rel="noopener noreferrer" href="https://vuejs.org">
          <img src="@/assets/svg/vue.svg" alt="Vue.js" />
        </a>
        and
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLneRaqXoYlV10RUsU6uP-6cicVmwr0g2O&index=400"
        >
          <img src="@/assets/svg/heart.svg" alt="love" />
        </a>
      </p>
    </div>
    <div class="mentions">
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
  opacity: 87%;
}

.container-m {
  padding-top: 10px;
  margin-top: 30px;
  padding-bottom: 50px;
  color: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  opacity: 87%;
}

.contacts {
  display: flex;
  justify-content: space-around;
  align-items: center;
  vertical-align: middle;
  width: 100%;
}

.mentions p {
  font-size: 0.6em;
}
</style>
