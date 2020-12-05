<template>
  <span>
    <b-form-group>
      <b-form-textarea
        class="mt-1"
        v-model="form.title"
        autocomplete="off"
        rows="1"
        :state="form.title ? properTitleLength : null"
        max-rows="8"
        no-resize
        required
        placeholder="* Catchy title..."
      ></b-form-textarea>
      <small
        v-if="form.title"
        :style="`color: ${properTitleLength ? 'green' : 'red'}`"
        >{{ titleLength }}/{{ maxTitleLength }}
      </small>
    </b-form-group>
    <b-form-group>
      <b-form-textarea
        id="textarea-auto-height"
        v-model="form.content"
        placeholder="* Cool content..."
        reqired
        :state="form.content ? properContentLength : null"
        rows="4"
        max-rows="50"
      ></b-form-textarea>
      <small
        v-if="form.content"
        :style="`color: ${properContentLength ? 'green' : 'red'}`"
        >{{ contentLength }}/{{ maxContentLength }}
      </small>
    </b-form-group>
    <b-form-tags
      autocomplete="off"
      remove-on-delete
      v-model="form.categories"
      tag-variant="dark"
      :placeholder="`Lowercase, ${minTagLength}-${maxTagLength} symbols`"
      :tag-validator="tagValidator"
      @tag-state="onTagState"
    ></b-form-tags>
  </span>
</template>
<script>
export default {
  props: { form: Object },
  watch: {
    validForm: function(val) {
      this.$emit("valid-form", val);
    },
  },
  computed: {
    validForm() {
      return (
        this.properTitleLength &&
        this.properContentLength &&
        !this.invalidTags.length &&
        !this.duplicateTags.length
      );
    },
    titleLength() {
      return this.form.title.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    contentLength() {
      return this.form.content.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properTitleLength() {
      return (
        this.titleLength >= this.minTitleLength &&
        this.titleLength <= this.maxTitleLength
      );
    },
    properContentLength() {
      return (
        this.contentLength >= this.minContentLength &&
        this.contentLength <= this.maxContentLength
      );
    },
  },
  data() {
    return {
      invalidTags: [],
      duplicateTags: [],

      minTitleLength: 5,
      maxTitleLength: 300,

      minContentLength: 5,
      maxContentLength: 2000,

      minTagLength: 3,
      maxTagLength: 20,
    };
  },
  created() {
    this.$emit("valid-form", this.validForm);
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
  },
};
</script>
