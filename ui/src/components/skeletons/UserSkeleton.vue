<template>
  <div class="grid">
    <div class="columns">
      <div class="main-col">
        <div>
          <b-skeleton
            variant="dark"
            class="mt-2 ml-2"
            type="button"
            width="300px"
            height="45px"
            animation="none"
          ></b-skeleton>
          <div v-for="index in 5" v-bind:key="index">
            <div :class="`${isMobile() ? 'card-m' : 'card'}`">
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
        <!-- <div v-else>
          <b-skeleton-img></b-skeleton-img>
        </div> -->
      </div>
      <div class="info-col">
        <div :class="`${isMobile() ? 'card-m' : 'card'}`">
          <b-skeleton
            variant="dark"
            type="avatar"
            style="display: block; margin: 5px auto 10px auto"
            width="85px"
            height="85px"
          ></b-skeleton>
          <b-skeleton
            style="display: block; margin: 5px auto 10px auto"
            variant="dark"
            animation="wave"
            width="45%"
          ></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="85%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="65%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="70%"></b-skeleton>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "UserSkeleton",
};
</script>
