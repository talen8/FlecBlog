<script lang="ts" setup>
import { parseJSON } from '@/utils/json';

const { blogConfig } = useSysConfig();

const allContacts = computed(() => {
  const socialList = parseJSON<
    Array<{ name: string; url: string; icon: string; position?: string }>
  >(blogConfig.value.footer_social, []);
  return socialList.filter(item => item.url && item.url.trim() !== '');
});

const contactsLeft = computed(() => {
  return allContacts.value.filter(item => item.position === 'left' || item.position === '左');
});

const contactsRight = computed(() => {
  return allContacts.value.filter(item => item.position === 'right' || item.position === '右');
});

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};
</script>

<template>
  <div class="footer-contact">
    <div v-for="contact in contactsLeft" :key="contact.name" class="contact_link">
      <a
        :href="contact.url"
        :aria-label="`访问 ${contact.name}`"
        target="_blank"
        rel="noopener noreferrer"
      >
        <i :class="'ri-' + contact.icon" aria-hidden="true" />
      </a>
    </div>
    <button class="footer_logo" aria-label="返回顶部" type="button" @click="scrollToTop">
      <svg
        id="Layer_1"
        viewBox="0 0 512 512"
        version="1.1"
        xml:space="preserve"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="Start_Up">
          <!-- Background Circle -->
          <circle id="Background_21_" cx="254.8" cy="255.7" r="215.6" style="fill: #eeedf2" />

          <!-- Clouds Group -->
          <g id="Clouds_4_">
            <!-- Left Cloud (Animated) -->
            <g>
              <animateTransform
                attributeName="transform"
                type="translate"
                values="0 0; -20 0; 0 0"
                dur="15s"
                repeatCount="indefinite"
                begin="0s"
              />
              <!-- Cloud Fill -->
              <path
                d="M195.9 142.7c0 11.8-9.6 21.4-21.4 21.4H86.7c-9.4 0-17-7.6-17-17 0-.6 0-1.1.1-1.7 0-.4.1-.8.2-1.2 1.2-6.8 6.4-12.3 13.1-13.7 1.2-.3 2.4-.4 3.6-.4h.2c-.1-.9-.2-1.9-.2-2.9 0-1.6.2-3.1.4-4.6 2.1-10.8 11.6-18.9 23-18.9 1.5 0 3 .1 4.4.4 3.3-10 11-17.9 20.9-21.5 3.1-1.1 6.4-1.8 9.8-1.9h1.9c18.6.1 33.7 15.3 33.7 34 0 2.5-.3 4.9-.8 7.2 9.2 2.6 15.9 10.9 15.9 20.8z"
                style="fill: #d3d5e7"
              />
              <!-- Cloud Outline 1 -->
              <path
                d="M180.5 132.4c0 10.8-8.8 19.6-19.6 19.6H80.8c-8.6 0-15.5-6.9-15.5-15.5 0-.5 0-1 .1-1.6 0-.4.1-.7.1-1.1 1.1-6.2 5.8-11.2 12-12.5 1.1-.2 2.2-.4 3.3-.4h.2c-.1-.9-.2-1.7-.2-2.6 0-1.4.1-2.8.4-4.2 1.9-9.9 10.6-17.3 21-17.3 1.4 0 2.7.1 4 .4 3-9.1 10.1-16.4 19.1-19.6 2.8-1 5.8-1.6 8.9-1.8h1.7c17 .1 30.8 14 30.8 31 0 2.3-.2 4.5-.7 6.6 8.4 2.3 14.5 9.9 14.5 19z"
                style="
                  fill: #fff;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <!-- Cloud Outline 2 -->
              <path
                d="M166.7 106.4c0 8.9-3.7 16.9-9.7 22.6"
                style="
                  fill: #fff;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <!-- Cloud Detail Line -->
              <path
                d="M101.6 96.7c9.8 0 18.1 6.8 20.3 15.9"
                style="
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
            </g>

            <!-- Right Cloud (Animated) -->
            <g>
              <animateTransform
                attributeName="transform"
                type="translate"
                values="0 0; 20 0; 0 0"
                dur="18s"
                repeatCount="indefinite"
                begin="3s"
              />
              <!-- Cloud Fill -->
              <path
                d="M473.3 255.9c0 11.8-9.6 21.4-21.4 21.4h-87.7c-9.4 0-17-7.6-17-17 0-.6 0-1.1.1-1.7 0-.4.1-.8.2-1.2 1.2-6.8 6.4-12.3 13.1-13.7 1.2-.3 2.4-.4 3.6-.4h.2c-.1-.9-.2-1.9-.2-2.9 0-1.6.2-3.1.4-4.6 2.1-10.8 11.6-18.9 23-18.9 1.5 0 3 .1 4.4.4 3.3-10 11-17.9 20.9-21.5 3.1-1.1 6.4-1.8 9.8-1.9h1.9c18.6.1 33.7 15.3 33.7 34 0 2.5-.3 4.9-.8 7.2 9.1 2.6 15.8 10.9 15.8 20.8z"
                style="fill: #d3d5e7"
              />
              <!-- Cloud Outline 1 -->
              <path
                d="M457.2 246.1c0 10.8-8.8 19.6-19.6 19.6h-80.1c-8.6 0-15.5-6.9-15.5-15.5 0-.5 0-1 .1-1.6 0-.4.1-.7.1-1.1 1.1-6.2 5.8-11.2 12-12.5 1.1-.2 2.2-.4 3.3-.4h.2c-.1-.9-.2-1.7-.2-2.6 0-1.4.1-2.8.4-4.2 1.9-9.9 10.6-17.3 21-17.3 1.4 0 2.7.1 4 .4 3-9.1 10.1-16.4 19.1-19.6 2.8-1 5.8-1.6 8.9-1.8h1.7c17 .1 30.8 14 30.8 31 0 2.3-.2 4.5-.7 6.6 8.4 2.4 14.5 10 14.5 19z"
                style="
                  fill: #fff;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <!-- Cloud Outline 2 -->
              <path
                d="M443.4 220.2c0 8.9-3.7 16.9-9.7 22.6"
                style="
                  fill: #fff;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <!-- Cloud Detail Line -->
              <path
                d="M378.3 210.5c9.8 0 18.1 6.8 20.3 15.9"
                style="
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
            </g>
          </g>

          <!-- Rocket Group (Animated Vertical Bounce) -->
          <g id="Rocket">
            <animateTransform
              xlink:href="#Rocket"
              attributeName="transform"
              type="translate"
              values="0 0; 0 -30; 0 30; 0 0"
              dur="3s"
              repeatCount="indefinite"
              begin="0s"
            />

            <!-- Rocket Sound Effects -->
            <g id="Sound_2_">
              <path
                id="Left_56_"
                d="m224.4 424.8-11.5 13.6c-.1.1 0 .4.2.4l15.7 1.4c.1 0 .1.1 0 .1l-11.5 10.8"
                style="
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <path
                id="Middle_60_"
                d="m254.6 418.1 8.1 15.8c.1.2-.1.4-.2.3l-15.6-2.2c-.1 0-.1.1-.1.1l8.7 13.1"
                style="
                  opacity: 0.2;
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <path
                id="Right_54_"
                d="m282.1 438 8.1 15.8c.1.2-.1.4-.2.3l-15.6-2.2c-.1 0-.1.1-.1.1l8.7 13.1"
                style="
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
            </g>

            <!-- Rocket Shadow -->
            <path
              id="Shadow_62_"
              d="M383.4 367.3h-50.9s-3.8 12.5-17.7 11.5v14.6c0 5.7-4.6 10.3-10.3 10.3H230c-5.7 0-10.3-4.6-10.3-10.3v-14.6c-13.9 1-17.7-11.5-17.7-11.5h-50.9c-1.1 0-2-1-1.8-2.1l3.4-22.1c1.3-8.6 5-16.6 10.5-23.3l34.1-40.5c3-3.6 4.7-8.1 4.7-12.8V211c0-23 6.1-45.6 17.7-65.4l24.3-41.7c10.4-17.8 36.1-17.8 46.5 0l24.3 41.7c11.6 19.9 17.7 42.4 17.7 65.4v55.5c0 4.7 1.6 9.2 4.7 12.8l34.1 40.5c5.6 6.6 9.2 14.7 10.5 23.3l3.4 22.1c.2 1.1-.7 2.1-1.8 2.1z"
              style="fill: #d3d5e7"
            />

            <!-- Rocket Launch Tube -->
            <g id="Launch_tube">
              <path
                id="Shape_90_"
                d="M302.5 380v14.6c0 5.7-4.6 10.3-10.3 10.3h-74.5c-5.7 0-10.3-4.6-10.3-10.3V380h95.1z"
                style="
                  fill: #da867d;
                  stroke: #1e247e;
                  stroke-width: 6;
                  stroke-linecap: round;
                  stroke-linejoin: round;
                  stroke-miterlimit: 10;
                "
              />
              <g id="Lines_38_">
                <path
                  id="Left_55_"
                  d="M208.5 392.5h13.4"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 2.5;
                    stroke-linecap: round;
                    stroke-miterlimit: 10;
                  "
                />
                <path
                  id="Right_53_"
                  d="M289.3 392.5h13.5"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 2.5;
                    stroke-linecap: round;
                    stroke-miterlimit: 10;
                  "
                />
              </g>
            </g>

            <!-- Rocket Body -->
            <g id="Body_13_">
              <path
                id="Shape_89_"
                d="M371.1 368.5h-50.9s-3.8 12.5-17.7 11.5h-95.1c-13.9 1-17.7-11.5-17.7-11.5h-50.9c-1.1 0-2-1-1.8-2.1l3.4-22.1c1.3-8.6 4.9-16.6 10.5-23.3l34.1-40.5c3-3.6 4.7-8.1 4.7-12.8v-55.5c0-23 6.1-45.6 17.7-65.4l24.3-41.7c10.4-17.8 36.1-17.8 46.5 0l24.3 41.7c11.6 19.9 17.7 42.4 17.7 65.4v55.5c0 4.7 1.6 9.2 4.7 12.8L359 321c5.6 6.6 9.2 14.7 10.5 23.3l3.4 22.1c.2 1.1-.7 2.1-1.8 2.1z"
                style="
                  fill: #fff;
                  stroke: #1e247e;
                  stroke-width: 6;
                  stroke-linecap: round;
                  stroke-linejoin: round;
                  stroke-miterlimit: 10;
                "
              />

              <!-- Rocket Windows (Animated Opacity) -->
              <g id="Windows_1_">
                <g id="Bottom_73_">
                  <circle
                    id="Shape_88_"
                    cx="254.7"
                    cy="281.6"
                    r="20.5"
                    style="
                      fill: #69aef8;
                      stroke: #1e247e;
                      stroke-width: 6;
                      stroke-linejoin: round;
                      stroke-miterlimit: 10;
                    "
                  />
                  <animate
                    xlink:href="#Shape_88_"
                    attributeName="opacity"
                    values="1; 0.7; 1"
                    dur="2s"
                    repeatCount="indefinite"
                  />
                </g>
                <g id="Top_75_">
                  <circle
                    id="Shape_87_"
                    cx="254.9"
                    cy="215.8"
                    r="20.5"
                    style="
                      fill: #69aef8;
                      stroke: #1e247e;
                      stroke-width: 6;
                      stroke-linejoin: round;
                      stroke-miterlimit: 10;
                    "
                  />
                  <animate
                    xlink:href="#Shape_87_"
                    attributeName="opacity"
                    values="1; 0.7; 1"
                    dur="3s"
                    repeatCount="indefinite"
                  />
                </g>
              </g>

              <!-- Rocket Body Lines -->
              <g id="Lines_37_">
                <path
                  d="m213.6 337.5 18.5 14.2M213.6 346.4l9.2 7M227 338.7l5.1 3.9"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 2.5;
                    stroke-linecap: round;
                    stroke-miterlimit: 10;
                  "
                />
              </g>

              <!-- Rocket Vertical Line -->
              <g id="Vertical_line_6_">
                <path
                  id="Bottom_72_"
                  d="M308.9 243.5v5.4"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 2.5;
                    stroke-linecap: round;
                    stroke-miterlimit: 10;
                  "
                />
                <path
                  id="Top_74_"
                  d="M308.9 221.8v17.3"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 2.5;
                    stroke-linecap: round;
                    stroke-miterlimit: 10;
                  "
                />
              </g>

              <!-- Rocket Top Line -->
              <path
                id="Top_line_6_"
                d="M260.7 101.2c2.6.5 5.2 1.7 7.1 3.8"
                style="
                  fill: none;
                  stroke: #1e247e;
                  stroke-width: 2.5;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />

              <!-- Rocket Fins -->
              <g id="Fins">
                <path
                  id="Left_54_"
                  d="M189.7 368.5V332"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 6;
                    stroke-linecap: round;
                    stroke-linejoin: round;
                    stroke-miterlimit: 10;
                  "
                />
                <path
                  id="Right__x23_2"
                  d="M319.7 368.5V332"
                  style="
                    fill: none;
                    stroke: #1e247e;
                    stroke-width: 6;
                    stroke-linecap: round;
                    stroke-linejoin: round;
                    stroke-miterlimit: 10;
                  "
                />
                <path
                  id="Right__x23_1"
                  d="M319.7 321.5c3.9 0 3.9-6 0-6s-3.9 6 0 6z"
                  style="fill: #1e247e"
                />
              </g>
            </g>

            <!-- Rocket Sparkles (Animated Opacity) -->
            <g id="Sparkles_22_">
              <path
                id="_x33_rd_23_"
                d="m314.7 83.5-8.4 2.5"
                style="
                  fill: none;
                  stroke: #fff;
                  stroke-width: 7;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <animate
                xlink:href="#_x33_rd_23_"
                attributeName="opacity"
                values="1; 0; 1"
                dur="1.5s"
                repeatCount="indefinite"
                begin="0s"
              />
              <path
                id="_x32_nd_24_"
                d="m305.7 71.5-6.7 7.1"
                style="
                  fill: none;
                  stroke: #fff;
                  stroke-width: 7;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <animate
                xlink:href="#_x32_nd_24_"
                attributeName="opacity"
                values="1; 0; 1"
                dur="2s"
                repeatCount="indefinite"
                begin="0.5s"
              />
              <path
                id="_x31_st_24_"
                d="m290.4 72.5 1.4-9.2"
                style="
                  fill: none;
                  stroke: #fff;
                  stroke-width: 7;
                  stroke-linecap: round;
                  stroke-miterlimit: 10;
                "
              />
              <animate
                xlink:href="#_x31_st_24_"
                attributeName="opacity"
                values="1; 0; 1"
                dur="1.8s"
                repeatCount="indefinite"
                begin="1s"
              />
            </g>
          </g>
        </g>
      </svg>
    </button>
    <div v-for="contact in contactsRight" :key="contact.name" class="contact_link">
      <a
        :href="contact.url"
        :aria-label="`访问 ${contact.name}`"
        target="_blank"
        rel="noopener noreferrer"
      >
        <i :class="'ri-' + contact.icon" aria-hidden="true" />
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.footer-contact {
  justify-content: center;
  display: flex;
  align-items: center;

  .contact_link {
    display: flex;
    margin: 14px 27px;
    color: var(--flec-footer-icon);
    border-radius: 20%;
    width: 38px;
    height: 38px;
    background: var(--flec-card-bg);
    justify-content: center;
    align-items: center;
    transition: 0.3s;
    cursor: pointer;

    &:hover {
      background: var(--flec-footer-icon-bg-hover);
      transform: scale(1.1);

      i {
        color: var(--flec-footer-icon-hover);
      }
    }

    a {
      display: flex;
      width: 100%;
      height: 100%;
      justify-content: center;
      align-items: center;
      color: inherit;
    }

    i {
      font-size: 1.4rem;
      pointer-events: none;
    }
  }

  .footer_logo {
    width: 60px;
    height: 60px;
    margin: 0 14px;
    cursor: pointer;
    transition: cubic-bezier(0, 0, 0, 1.29) 0.5s;
    border-radius: 50%;
    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    border: none;
    background: transparent;
    padding: 0;

    &:hover {
      transform: scale(1.1);
    }

    &:focus {
      outline: 2px solid var(--flec-btn);
      outline-offset: 2px;
    }
  }
}

// 响应式设置
@media screen and (max-width: 768px) {
  .footer-contact {
    flex-wrap: wrap;
    justify-content: center;
    gap: 60px;

    .contact_link {
      margin: 0;
    }

    .footer_logo {
      display: none;
    }
  }
}
</style>
