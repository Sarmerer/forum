<template>
  <b-overlay
    :show="requesting"
    rounded
    opacity="0.6"
    spinner-small
    variant="dark"
    spinner-variant="light"
  >
    <template #overlay>
      <div class="text-center">
        <b-icon icon="stopwatch" font-scale="3" animation="cylon"></b-icon>
        <p id="cancel-label">{{ postCreationState }}</p>
      </div>
    </template>
    <b-form @submit.prevent="submit">
      <div class="tabs">
        <b-button
          size="sm"
          class="tab text-white-87"
          :disabled="contentType === 'text'"
          :variant="contentType === 'text' ? 'outline-light' : 'dark'"
          @click="contentType = 'text'"
          ><b-icon-newspaper></b-icon-newspaper> Write</b-button
        >
        <b-button
          size="sm"
          class="tab text-white-87"
          :disabled="contentType === 'image'"
          :variant="contentType === 'image' ? 'outline-light' : 'dark'"
          @click="contentType = 'image'"
          ><b-icon-image></b-icon-image> Upload</b-button
        >
      </div>
      <b-form-group>
        <b-form-textarea
          id="textarea"
          aria-describedby="title-help-block"
          class="mt-1"
          v-model="form.title"
          autocomplete="off"
          rows="1"
          :state="form.title ? properTitleLength : null"
          max-rows="8"
          no-resize
          required
          placeholder="Catchy title"
        ></b-form-textarea>

        <b-form-text id="title-help-block" v-if="form.title">
          <span :style="`color: ${properTitleLength ? 'white' : 'red'}`"
            >{{ titleLength }}/{{ maxTitleLength }}
          </span>
        </b-form-text>
      </b-form-group>
      <div class="mb-3">
        <template v-if="contentType === 'text'">
          <b-form-group>
            <b-form-textarea
              aria-describedby="content-help-block"
              v-model="form.content"
              placeholder="Cool content (optional)"
              reqired
              :state="form.content ? properContentLength : null"
              rows="4"
              max-rows="50"
            ></b-form-textarea>
            <b-form-text id="content-help-block" v-show="form.content">
              <span :style="`color: ${properContentLength ? 'white' : 'red'}`"
                >{{ contentLength }}/{{ maxContentLength }}
              </span>
            </b-form-text>
          </b-form-group>
        </template>
        <template v-if="contentType === 'image'">
          <div
            class="image-upload-wrapper"
            @dragover.prevent
            @drop.prevent="onImageDrop"
          >
            <div v-if="imagePreview">
              <b-img-lazy
                :src="imagePreview"
                center
                rounded
                fluid-grow
              ></b-img-lazy>
              <b-button
                class="dismiss"
                @click="clearSelectedImage"
                variant="outline-light"
              >
                <b-icon-x></b-icon-x>
              </b-button>
            </div>
            <div v-if="!imagePreview" class="file-input">
              <label for="file">Drop an image or click to select</label>
              <input
                class="image-upload"
                type="file"
                id="file"
                accept="image/*"
                @change.prevent="onImageChange"
              />
            </div>
          </div>
        </template>
      </div>
      <b-form-tags
        aria-describedby="tags-help-block"
        autocomplete="off"
        separator=" "
        remove-on-delete
        v-model="form.categories"
        tag-variant="dark"
        placeholder="Tags (optional)"
        :tag-validator="tagValidator"
        @tag-state="onTagState"
      ></b-form-tags>
      <b-form-text id="tags-help-block">
        Lowercase, {{ minTagLength }}-{{ maxTagLength }} symbols
      </b-form-text>
      <slot name="buttons" :validForm="validForm"></slot>
    </b-form>
  </b-overlay>
</template>
<script>
import api from "@/api/api";

export default {
  props: { formData: Object, edit: Boolean },
  computed: {
    validForm() {
      return (
        this.properTitleLength &&
        this.properContentLength &&
        !this.invalidTags.length &&
        !this.duplicateTags.length &&
        this.validImage
      );
    },
    titleLength() {
      return this.form?.title.replace(/(\r\n|\n|\r|\s)/g, "").length || 0;
    },
    contentLength() {
      return this.form?.content.replace(/(\r\n|\n|\r|\s)/g, "").length || 0;
    },
    properTitleLength() {
      return (
        this.titleLength >= this.minTitleLength &&
        this.titleLength <= this.maxTitleLength
      );
    },
    properContentLength() {
      return this.contentLength <= this.maxContentLength;
    },
    validImage() {
      return this.contentType === "image" ? this.imagePreview ? true : this.selectedImage != null : true;
    },
  },
  data() {
    return {
      requesting: false,
      contentType: this.formData?.is_image ? "image" : "text" || "text",
      form: {
        id: this.formData?.id || 0,
        title: this.formData?.title || "",
        content: this.formData?.is_image ? "" : this.formData?.content || "",
        is_image: this.formData?.is_image || false,
        categories: this.formData?.categories?.map((c) => c.name) || [],
      },

      invalidTags: [],
      duplicateTags: [],

      minTitleLength: 5,
      maxTitleLength: 300,

      minContentLength: 5,
      maxContentLength: 2000,

      minTagLength: 3,
      maxTagLength: 20,

      selectedImage: null,
      imagePreview: this.formData?.is_image ? this.formData?.content : "" || "",
      imageChanged: false,

      postCreationState: "",
    };
  },

  methods: {
    onTagState(_valid, invalid, duplicate) {
      this.invalidTags = invalid;
      this.duplicateTags = duplicate;
    },
    tagValidator(tag) {
      return (
        tag === tag.toLowerCase() &&
        tag.length >= this.minTagLength &&
        tag.length <= this.maxTagLength
      );
    },
    onImageChange(event) {
      event.preventDefault();
      if (!event?.target?.files[0]) return;
      this.imageChanged = true;
      this.selectedImage = event?.target?.files[0];
      this.imagePreview = URL.createObjectURL(this.selectedImage);
    },
    onImageDrop(event) {
      const files = event?.dataTransfer?.files;
      if (!files[0]) return;
      this.imageChanged = true;
      this.selectedImage = files[0];
      this.imagePreview = URL.createObjectURL(this.selectedImage);
    },
    clearSelectedImage() {
      this.selectedImage = null;
      this.imagePreview = "";
    },
    submit() {
      if (this.requesting) return;
      let action = this.edit
        ? (action = this.updatePost)
        : this.createPostWithText;

      if (this.contentType === "image" && this.imageChanged) {
        this.uploadImage(action);
      } else {
        action();
      }
    },
    async uploadImage(callback) {
      if (this.requesting) return;
      this.requesting = true;
      const fd = new FormData();
      fd.append("image", this.selectedImage, this.selectedImage.name);
      return api
        .post("image/upload", fd, {
          onUploadProgress(event) {
            this.postCreationState = `Uploading image: ${Math.round(
              event.loaded / event.total
            ) * 100}%`;
          },
          headers: {
            "Content-Type": `multipart/form-data;boundary=${this.selectedImage._boundary}`,
          },
        })
        .then((response) => {
          this.form.content = response.data.data;
	  callback()
        })
        .catch(() => {
	this.$bvToast.toast("This image is too heavy!", {
            title: "Oops!",
            variant: "danger",
            solid: true,
          });
	}).then(this.requesting = false)
    },
    async createPostWithText() {
      if (this.requesting) return;
      this.requesting = true;
      this.postCreationState = "Creating post...";
      return api
        .post("post/create", {
          title: this.form.title,
          content: this.form.content,
          categories: this.form.categories,
          is_image: this.contentType === "image",
        })
        .then((response) => {
          this.$router.push({
            name: "Post",
            params: {
              id: response.data?.data?.id,
              postData: response?.data?.data,
            },
          });
        })
        .then((this.requesting = false));
    },

    async updatePost() {
      if (this.requesting) return;
      this.postCreationState = "Updating post...";
      this.requesting = true;
      return await api
        .put("post/update", {
          id: this.form.id,
          title: this.form.title,
          content: this.form.content,
          is_image: this.contentType === "image",
          categories: this.form.categories,
        })
        .then((response) => {
          if (response?.data?.data) {
            this.$emit("post-update", response.data.data);
          }
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast("You need to be logged in, to update posts!", {
              title: "Oops!",
              variant: "danger",
              solid: true,
            });
        })
        .then((this.requesting = false));
    },
  },
};
</script>
<style lang="scss" scoped>
.tabs {
  display: flex;
  justify-content: space-around;
  gap: 5px;
}

.tab {
  width: 100%;
}

.image-upload-wrapper {
  width: 100%;
  color: #fff;
  padding: 5px;
  text-align: center;
  border-radius: 4px;
  border: 2px dashed rgba(255, 255, 255, 0.5);
  position: relative;
}

.image-upload {
  display: none;
}

.file-input {
  width: 100%;
  margin: auto;
  height: 200px;
  position: relative;
  label,
  input {
    color: white;
    opacity: 0.87;
    width: 100%;
    margin-top: 80px;
    cursor: pointer;
  }

  input {
    opacity: 0;
    z-index: -2;
  }
}

.dismiss {
  position: absolute;
  top: 0;
  right: 0;
}
</style>
