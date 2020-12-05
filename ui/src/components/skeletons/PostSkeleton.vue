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
              <b-skeleton
                variant="dark"
                animation="wave"
                width="45%"
              ></b-skeleton
            ></b-col>
          </b-row>
          <b-skeleton variant="dark" animation="wave" width="100%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="100%"></b-skeleton>
        </div>
      </div>
      <div class="main-col">
        <div :class="`${isMobile() ? 'card-m' : 'card'}`">
          <b-skeleton variant="dark" animation="wave" width="75%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="65%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
        </div>
        <div :class="`${isMobile() ? 'card-m' : 'card'}`">
          <b-skeleton
            type="input"
            variant="dark"
            animation="wave"
            width="100%"
            height="50px"
          ></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
          <b-row class="ml-1">
            <b-skeleton
              type="avatar"
              variant="dark"
              animation="wave"
              width="15px"
              height="15px"
              class="mr-1"
            ></b-skeleton>
            <b-skeleton
              variant="dark"
              animation="wave"
              width="15%"
            ></b-skeleton>
          </b-row>
          <b-skeleton variant="dark" animation="wave" width="25%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
          <b-row class="ml-1">
            <b-skeleton
              type="avatar"
              variant="dark"
              animation="wave"
              width="15px"
              height="15px"
              class="mr-1"
            ></b-skeleton>
            <b-skeleton
              variant="dark"
              animation="wave"
              width="25%"
            ></b-skeleton>
          </b-row>
          <b-skeleton variant="dark" animation="wave" width="20%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
          <b-row class="ml-1">
            <b-skeleton
              type="avatar"
              variant="dark"
              animation="wave"
              width="15px"
              height="15px"
              class="mr-1"
            ></b-skeleton>
            <b-skeleton
              variant="dark"
              animation="wave"
              width="10%"
            ></b-skeleton>
          </b-row>
          <b-skeleton variant="dark" animation="wave" width="30%"></b-skeleton>
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
  name: "PostSkeleton",
};
</script>
