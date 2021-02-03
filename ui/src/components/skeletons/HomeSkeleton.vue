<template>
  <div class="grid">
    <div class="columns">
      <div class="main-col">
        <div>
          <div
            :align="isMobile() ? 'center' : 'right'"
            :class="{ 'mr-4': !isMobile() }"
          >
            <b-skeleton
              type="button"
              variant="dark"
              animation="none"
              width="200px"
            ></b-skeleton>
          </div>

          <div v-for="index in 5" :key="'p' + index">
            <div :class="{ 'card-m': isMobile(), card: !isMobile() }">
              <b-skeleton
                variant="dark"
                animation="wave"
                width="85%"
              ></b-skeleton>
              <b-skeleton
                variant="dark"
                animation="wave"
                width="55%"
              ></b-skeleton>
              <b-skeleton
                variant="dark"
                animation="wave"
                width="70%"
              ></b-skeleton>
            </div>
          </div>
        </div>
      </div>
      <div class="info-col">
        <div :class="{ 'card-m': isMobile(), card: !isMobile() }">
          <h3 class="primary">RECENT</h3>
          <div v-for="index in 3" :key="'r' + index">
            <b-skeleton
              variant="dark"
              animation="wave"
              width="85%"
            ></b-skeleton>
            <b-skeleton
              variant="dark"
              animation="wave"
              width="55%"
            ></b-skeleton>
            <b-skeleton
              variant="dark"
              animation="wave"
              width="70%"
            ></b-skeleton>
            <br />
          </div>
        </div>
        <div :class="{ 'card-m': isMobile(), card: !isMobile() }">
          <h3 class="primary">CATEGORIES</h3>
          <b-skeleton variant="dark" animation="wave" width="85%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="70%"></b-skeleton>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "HomeSkeleton",
};
</script>
