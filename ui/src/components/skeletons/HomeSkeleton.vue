      <template>
  <div class="grid">
    <div class="columns">
      <div class="main-col">
        <div v-if="postsLength > 0">
          <b-skeleton
            type="button"
            variant="dark"
            style="display: block; margin: 15px 0 0 auto"
            animation="none"
            width="200px"
          ></b-skeleton>
          <div v-for="index in postsLength" v-bind:key="index">
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
          <h3 class="primary">RECENT</h3>
          <b-skeleton variant="dark" animation="wave" width="85%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="55%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="70%"></b-skeleton>
          <br />
          <b-skeleton variant="dark" animation="wave" width="80%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="45%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="90%"></b-skeleton>
        </div>
        <div :class="`${isMobile() ? 'card-m' : 'card'}`">
          <h3 class="primary">CATEGORIES</h3>
          <b-skeleton variant="dark" animation="wave" width="45%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="100%"></b-skeleton>
          <b-skeleton variant="dark" animation="wave" width="75%"></b-skeleton>
        </div>
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