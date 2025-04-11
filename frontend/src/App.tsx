import { useEffect, useState } from 'react';
import '@shopify/polaris/build/esm/styles.css';
import { AppProvider, Card, Page, Button, Text, Thumbnail, Tabs, BlockStack, InlineStack, Box, Divider } from '@shopify/polaris';
// import { MobileBackArrowMajor } from '@shopify/polaris-icons';
import translations from '@shopify/polaris/locales/en.json';
// import { authenticatedFetch } from "@shopify/app-bridge/utilities";
import {useAppBridge} from '@shopify/app-bridge-react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useLocation,
} from 'react-router-dom';

function AppPage() {
  const [selectedTab, setSelectedTab] = useState(0);
  const shopify = useAppBridge();

  const tabs = [
    {
      id: 'model-tab',
      content: 'Model',
    },
    {
      id: 'product-tab',
      content: 'Product',
    },
  ];

  const productImages = [
    'https://img.icons8.com/ios-filled/100/000000/glasses.png',
    'https://img.icons8.com/ios-filled/100/000000/eyeglasses.png',
    'https://img.icons8.com/ios-filled/100/000000/sunglasses.png',
  ];

  return (
    <Page title="VIRTUAL TRY-ON">
      <BlockStack gap="4">
        <Card padding="5">
          <BlockStack gap="4">
            <Box paddingBlockEnd="4">
              <Text variant="headingLg" as="h2" alignment="center">
                Try on glasses virtually
              </Text>
            </Box>
            
            <Box padding="4" background="bg-surface-secondary" borderRadius="3">
              <BlockStack gap="4" alignment="center">
                <div style={{ position: 'relative', maxWidth: '300px', margin: '0 auto' }}>
                  <img
                    src="https://images.unsplash.com/photo-1607746882042-944635dfe10e"
                    alt="Model wearing glasses"
                    style={{ width: '100%', height: 'auto', borderRadius: '12px', boxShadow: '0 4px 8px rgba(0,0,0,0.1)' }}
                  />
                  <div style={{
                    position: 'absolute', top: 16, left: 16, 
                    color: '#fff', fontWeight: 'bold',
                    background: 'rgba(0,0,0,0.5)', padding: '4px 8px',
                    borderRadius: '4px'
                  }}>Model</div>
                  <div style={{
                    position: 'absolute', top: 16, right: 16, 
                    color: '#fff', cursor: 'pointer',
                    background: 'rgba(0,0,0,0.5)', padding: '4px 8px',
                    borderRadius: '4px'
                  }}>Cancel</div>
                </div>
                
                <Button size="large" primary onClick={() => {
                  shopify.toast.show('Blog post template generated');

                  fetch('/api/test', {
                    method: 'POST',
                    headers: {
                      'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({}),
                  })
                }}>
                  Take photo
                </Button>
              </BlockStack>
            </Box>
            
            <Divider />
            
            <Tabs tabs={tabs} selected={selectedTab} onSelect={setSelectedTab} fitted>
              {selectedTab === 0 ? (
                <Box padding="4">
                  <BlockStack gap="4">
                    <Text variant="headingMd" as="h2" alignment="center">
                      Upload your photo
                    </Text>
                    <Text alignment="center" color="subdued">
                      Take a selfie or upload a photo to see how glasses look on you
                    </Text>
                    <Button fullWidth>Upload image</Button>
                  </BlockStack>
                </Box>
              ) : (
                <Box padding="4">
                  <BlockStack gap="4">
                    <Text variant="headingMd" as="h2" alignment="center">
                      Available Products
                    </Text>
                    
                    <InlineStack gap="4" wrap={false} align="center" distribute="center">
                      {productImages.map((src, index) => (
                        <Card key={index} padding="3">
                          <BlockStack gap="2" alignment="center">
                            <Thumbnail
                              source={src}
                              alt={`Product ${index + 1}`}
                              size="large"
                            />
                            <Text variant="bodySm">Style {index + 1}</Text>
                          </BlockStack>
                        </Card>
                      ))}
                    </InlineStack>
                    
                    <Box paddingBlockStart="2">
                      <Button plain>View all products</Button>
                    </Box>
                  </BlockStack>
                </Box>
              )}
            </Tabs>
          </BlockStack>
        </Card>
      </BlockStack>
    </Page>
  );
}

function ProductPage() {
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const productId = query.get('id');

  const [images, setImages] = useState<string[]>([]);

  useEffect(() => {
    if (productId) {
      fetch(`/api/product/${productId}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then(res => res.json())
        .then(data => {
          // Assuming data.images is an array of image URLs
          setImages(data.images);
        });
    }
  }, [productId]);

  return <div>Product</div>;
}

function App() {
  return (
    <AppProvider i18n={translations} features={{ newDesignLanguage: true }}>
      <Router>
        <Routes>
          <Route path="/" element={<AppPage />} />
          <Route path="/product" element={<ProductPage />} />
        </Routes>
      </Router>
    </AppProvider>
  );
}

export default App;
