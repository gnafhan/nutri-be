package seeders

import (
	"app/src/model"
	"app/src/utils"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedArticles(db *gorm.DB) {
	var count int64
	db.Model(&model.Article{}).Count(&count)

	if count > 0 {
		log.Println("✅ Articles already seeded, skipping...")
		return
	}

	// Fixed user ID (admin/author)
	userID, _ := uuid.Parse("7ea4f83c-64c4-4079-ad44-2c3d1e06be2d")

	// Create article categories first
	categories := []model.ArticleCategory{
		{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Nutrition",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Diet",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Exercise",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Wellness",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Healthy Living",
		},
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Fatalf("Failed to seed article categories: %v", err)
	}

	// Sample dummy images URLs
	dummyImages := []string{
		"https://placehold.co/800x400/png?text=Health+Article+Image",
		"https://placehold.co/800x400/png?text=Nutrition+Guide",
		"https://placehold.co/800x400/png?text=Diet+Tips",
		"https://placehold.co/800x400/png?text=Exercise+Routines",
		"https://placehold.co/800x400/png?text=Wellness+Strategies",
	}

	// Sample article content in markdown format
	nutritionContent := `# Understanding Balanced Nutrition

## Introduction
Balanced nutrition is essential for maintaining good health. It involves consuming a variety of foods in the right proportions.

## Key Elements of a Balanced Diet
- **Proteins**: Essential for building and repairing tissues
- **Carbohydrates**: The body's primary source of energy
- **Fats**: Important for hormone production and vitamin absorption
- **Vitamins and Minerals**: Necessary for various bodily functions
- **Water**: Crucial for hydration and digestion

## Creating a Balanced Meal
A balanced meal should include:
1. A source of lean protein (chicken, fish, legumes)
2. Complex carbohydrates (whole grains, starchy vegetables)
3. Healthy fats (avocados, nuts, olive oil)
4. Fruits and vegetables (at least half your plate)

> Remember, the key is moderation and variety!

## Benefits of Balanced Nutrition
- Maintains energy levels throughout the day
- Supports immune function
- Promotes healthy growth and development
- Reduces the risk of chronic diseases

*Consult with a registered dietitian for personalized advice.*`

	dietContent := `# Effective and Sustainable Diet Approaches

## Understanding Diets
Different diets work for different people. The best diet is one that you can maintain long-term.

## Popular Diet Types
### Mediterranean Diet
- Rich in fruits, vegetables, whole grains, and olive oil
- Moderate amounts of fish, poultry, and dairy
- Limited red meat consumption

### Plant-Based Diet
- Focuses on foods derived from plant sources
- Can include vegetarian or vegan approaches
- High in fiber and antioxidants

### DASH Diet
- Designed to help lower blood pressure
- Emphasizes fruits, vegetables, and low-fat dairy
- Limits sodium, added sugars, and red meat

## Creating Sustainable Habits
1. Make gradual changes to your eating patterns
2. Focus on adding healthy foods rather than restriction
3. Listen to your body's hunger and fullness cues
4. Allow flexibility for special occasions

> Sustainable changes lead to lasting results!

## Warning Signs of Unhealthy Diets
- Extreme restriction of food groups
- Rapid weight loss
- Obsessive food behaviors
- Fatigue or malnutrition symptoms

*Always consult healthcare professionals before starting any diet.*`

	exerciseContent := `# Building an Effective Exercise Routine

## Benefits of Regular Exercise
Regular physical activity improves cardiovascular health, builds strength, and enhances mood.

## Types of Exercise
### Cardiovascular Training
- Running, swimming, cycling
- Improves heart health and endurance
- Aim for 150 minutes of moderate activity weekly

### Strength Training
- Weight lifting, resistance bands, bodyweight exercises
- Builds muscle mass and bone density
- Target each major muscle group 2-3 times weekly

### Flexibility and Balance
- Yoga, Pilates, tai chi
- Improves range of motion and prevents injury
- Often overlooked but essential components

## Creating a Balanced Fitness Plan
1. Combine different exercise types throughout the week
2. Start slowly and progressively increase intensity
3. Include rest days for recovery
4. Set specific, achievable goals

> Consistency matters more than intensity!

## Sample Weekly Schedule
| Day | Activity |
|-----|----------|
| Monday | Strength Training |
| Tuesday | Cardio |
| Wednesday | Yoga/Flexibility |
| Thursday | Strength Training |
| Friday | Cardio |
| Weekend | Active Recovery |

*Remember to warm up before and cool down after exercise sessions.*`

	wellnessContent := `# Holistic Wellness Strategies

## What is Wellness?
Wellness encompasses physical, mental, emotional, and social well-being.

## Key Dimensions of Wellness
### Physical Wellness
- Regular exercise and proper nutrition
- Adequate sleep and rest
- Preventive healthcare

### Mental and Emotional Wellness
- Stress management techniques
- Mindfulness and meditation
- Seeking support when needed

### Social Wellness
- Building meaningful relationships
- Community involvement
- Work-life balance

## Daily Wellness Practices
1. Morning meditation or gratitude practice
2. Regular movement breaks throughout the day
3. Digital detox periods
4. Quality sleep routine

> Small daily habits create profound long-term benefits.

## Signs Your Wellness Needs Attention
- Persistent fatigue or low energy
- Trouble sleeping or concentrating
- Feeling disconnected or isolated
- Increased irritability or mood swings

*Wellness is a journey, not a destination. Be patient with yourself.*`

	healthyLivingContent := `# Healthy Living for Longevity

## Foundations of Healthy Living
Healthy living involves making conscious choices that promote overall well-being and longevity.

## Key Components
### Nutritious Eating Patterns
- Focus on whole, minimally processed foods
- Regular, balanced meals
- Mindful eating practices

### Active Lifestyle
- Regular movement throughout the day
- Finding enjoyable physical activities
- Avoiding prolonged sedentary periods

### Quality Sleep
- Consistent sleep schedule
- Creating a restful environment
- Managing sleep disruptors

### Stress Management
- Identifying stress triggers
- Developing healthy coping mechanisms
- Setting appropriate boundaries

## Practical Health-Promoting Habits
1. Stay hydrated throughout the day
2. Practice good posture
3. Take regular breaks from screens
4. Spend time outdoors daily

> Health is wealth - invest in it daily!

## Creating Sustainable Changes
- Focus on one small change at a time
- Build a supportive environment
- Track progress but avoid perfectionism
- Celebrate small victories

*Remember that health looks different for everyone. Focus on what works for you.*`

	// Collection of content pieces
	contentOptions := []string{
		nutritionContent,
		dietContent,
		exerciseContent,
		wellnessContent,
		healthyLivingContent,
	}

	// Article titles
	titles := []string{
		"The Complete Guide to Balanced Nutrition",
		"Understanding Macronutrients and Micronutrients",
		"How to Create a Sustainable Diet Plan",
		"The Truth About Popular Diet Trends",
		"Building an Effective Exercise Routine",
		"Strength Training for Beginners",
		"Mindfulness and Mental Wellness",
		"Sleep Quality and Your Health",
		"Hydration: The Foundation of Good Health",
		"Stress Management Techniques That Actually Work",
		"Healthy Eating on a Budget",
		"Plant-Based Nutrition Guide",
		"The Connection Between Gut Health and Overall Wellness",
		"Finding the Right Exercise for Your Body Type",
		"Healthy Aging: Nutrition and Lifestyle Tips",
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create articles
	var articles []model.Article
	now := time.Now()

	for i, title := range titles {
		// Select random category
		categoryIndex := rand.Intn(len(categories))
		categoryID := categories[categoryIndex].ID

		// Select image and content
		imageIndex := rand.Intn(len(dummyImages))
		contentIndex := rand.Intn(len(contentOptions))
		imageURL := dummyImages[imageIndex]
		content := contentOptions[contentIndex]

		// Generate publish date (between 1-30 days ago)
		daysAgo := rand.Intn(30) + 1
		publishedAt := now.Add(-time.Hour * 24 * time.Duration(daysAgo))

		// Create slug from title
		slug := utils.Slugify(title)

		article := model.Article{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       title,
			CategoryID:  &categoryID,
			Slug:        slug,
			Image:       &imageURL,
			Content:     content,
			PublishedAt: &publishedAt,
			CreatedAt:   publishedAt,
			UpdatedAt:   publishedAt,
		}

		articles = append(articles, article)

		// Add a duplicate article with different title if needed to ensure we have enough
		if i < 5 {
			modifiedTitle := "Guide: " + title
			article := model.Article{
				ID:          uuid.New(),
				UserID:      userID,
				Title:       modifiedTitle,
				CategoryID:  &categoryID,
				Slug:        utils.Slugify(modifiedTitle),
				Image:       &imageURL,
				Content:     content,
				PublishedAt: &publishedAt,
				CreatedAt:   publishedAt,
				UpdatedAt:   publishedAt,
			}
			articles = append(articles, article)
		}
	}

	if err := db.Create(&articles).Error; err != nil {
		log.Fatalf("Failed to seed articles: %v", err)
	}

	log.Printf("✅ %d article categories and %d articles seeded successfully", len(categories), len(articles))
}
