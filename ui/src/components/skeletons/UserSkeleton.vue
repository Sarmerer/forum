<template>
  <div class="grid">
    <div class="columns">
      <div class="info-col" v-if="isMobile()">
        <div class="card-m">
          <b-row align-v="end" class="mb-3">
            <b-col cols="3">
              <b-skeleton
                variant="dark"
                type="avatar"
                width="55px"
                height="55px"
              ></b-skeleton
            ></b-col>
            <b-col>
              <b-skeleton
                variant="dark"
                animation="wave"
                width="50%"
              ></b-skeleton>
              <b-skeleton variant="dark" animation="wave" width="45%">
              </b-skeleton>
            </b-col>
          </b-row>
          <b-skeleton variant="dark" animation="wave" width="100%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="100%"></b-skeleton>
        </div>
      </div>
      <div class="main-col">
        <div>
          <div>
            <div :class="{ 'card-m': isMobile(), card: !isMobile() }">
              <b-skeleton
                variant="dark"
                type="button"
                animation="wave"
                width="100%"
              ></b-skeleton>
            </div>
          </div>
          <div v-for="index in 5" v-bind:key="index">
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
      <div class="info-col" v-if="!isMobile()">
        <div class="card">
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
