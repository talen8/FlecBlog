import { defineCollection, z } from '@nuxt/content'

const variantEnum = z.enum(['solid', 'outline', 'subtle', 'soft', 'ghost', 'link'])
const colorEnum = z.enum(['primary', 'secondary', 'neutral', 'error', 'warning', 'success', 'info'])
const sizeEnum = z.enum(['xs', 'sm', 'md', 'lg', 'xl'])
const orientationEnum = z.enum(['vertical', 'horizontal'])

const createBaseSchema = () => z.object({
  title: z.string().nonempty(),
  description: z.string().nonempty()
})

const createFeatureItemSchema = () => createBaseSchema().extend({
  icon: z.string().nonempty().editor({ input: 'icon' })
})

const createLinkSchema = () => z.object({
  label: z.string().nonempty(),
  to: z.string().nonempty(),
  icon: z.string().optional().editor({ input: 'icon' }),
  size: sizeEnum.optional(),
  trailing: z.boolean().optional(),
  target: z.string().optional(),
  color: colorEnum.optional(),
  variant: variantEnum.optional()
})

const createImageSchema = () => z.object({
  src: z.string().nonempty().editor({ input: 'media' }),
  alt: z.string().optional(),
  loading: z.enum(['lazy', 'eager']).optional(),
  srcset: z.string().optional()
})

export const collections = {
  index: defineCollection({
    source: '0.index.yml',
    type: 'page',
    schema: z.object({
      hero: z.object(({
        headline: z.string().optional(),
        command: z.string().optional(),
        links: z.array(createLinkSchema())
      })),
      sections: z.array(
        createBaseSchema().extend({
          id: z.string().nonempty(),
          orientation: orientationEnum.optional(),
          reverse: z.boolean().optional(),
          image: z.string().optional(),
          features: z.array(createFeatureItemSchema())
        })
      ),
      features: createBaseSchema().extend({
        items: z.array(createFeatureItemSchema())
      }),
      testimonials: createBaseSchema().extend({
        headline: z.string().optional(),
        items: z.array(
          z.object({
            quote: z.string().nonempty(),
            user: z.object({
              name: z.string().nonempty(),
              description: z.string().nonempty(),
              to: z.string().nonempty(),
              target: z.string().nonempty(),
              avatar: createImageSchema()
            })
          })
        )
      }),
      cta: createBaseSchema().extend({
        links: z.array(createLinkSchema())
      })
    })
  }),
  docs: defineCollection({
    source: '1.docs/**/*',
    type: 'page'
  }),
  pricing: defineCollection({
    source: '4.pricing.yml',
    type: 'page',
    schema: z.object({
      plans: z.array(
        z.object({
          title: z.string().nonempty(),
          description: z.string().nonempty(),
          price: z.object({
            month: z.string().nonempty(),
            year: z.string().nonempty()
          }),
          billing_period: z.string().nonempty(),
          billing_cycle: z.string().nonempty(),
          billingCycle: z.string().optional(),
          discount: z.string().optional(),
          button: createLinkSchema(),
          features: z.array(z.string().nonempty()),
          highlight: z.boolean().optional(),
          scale: z.boolean().optional()
        })
      ),
      logos: z.object({
        title: z.string().nonempty(),
        icons: z.array(z.string())
      }),
      faq: createBaseSchema().extend({
        items: z.array(
          z.object({
            label: z.string().nonempty(),
            content: z.string().nonempty()
          })
        )
      })
    })
  }),
  themes: defineCollection({
    source: '2.themes.yml',
    type: 'page',
    schema: z.object({
      items: z.array(
        z.object({
          name: z.string().nonempty(),
          display_name: z.string().nonempty(),
          slug: z.string().nonempty(),
          author: z.string().nonempty(),
          description: z.string().optional(),
          image: z.string().nonempty(),
          links: z.object({
            preview: z.string().optional(),
            download: z.string().optional(),
            source: z.string().optional()
          }).optional()
        })
      )
    })
  }),
  community: defineCollection({
    source: '5.community.yml',
    type: 'page',
    schema: z.object({
      items: z.array(
        z.object({
          title: z.string().nonempty(),
          description: z.string().nonempty(),
          icon: z.string().nonempty().editor({ input: 'icon' }),
          link: z.string().nonempty(),
          link_label: z.string().nonempty()
        })
      )
    })
  }),
  changelog: defineCollection({
    source: '3.changelog.yml',
    type: 'page'
  }),
  versions: defineCollection({
    source: '3.changelog/**/*',
    type: 'page',
    schema: z.object({
      title: z.string().nonempty(),
      description: z.string(),
      date: z.date(),
      image: z.string()
    })
  })
}
