      <template>
  <div class="grid">
    <div class="columns">
      <div class="main-col">
        <div v-if="postsLength > 0">
          <div v-for="index in postsLength" v-bind:key="index">
            <b-card>
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
            </b-card>
          </div>
        </div>
        <!-- <div v-else>
          <b-skeleton-img></b-skeleton-img>
        </div> -->
      </div>
      <div class="info-col">
        <b-card>
          <b-skeleton variant="dark" animation="wave" width="85%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="70%"></b-skeleton>
        </b-card>
        <b-card>
          <b-skeleton variant="dark" animation="wave" width="85%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="70%"></b-skeleton>
        </b-card>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "HomeSkeleton",
  props: ["postsLength"],
};
</script>