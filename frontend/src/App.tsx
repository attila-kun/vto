import { useState } from 'react'
import '@shopify/polaris/build/esm/styles.css';
import { AppProvider, Page, Card, TextField } from '@shopify/polaris'
import translations from "@shopify/polaris/locales/en.json";

function AppPage() {
  const [value, setValue] = useState('')

  return (
    <Page title="Polaris Skeleton App">
      <Card sectioned>
        <TextField
          label="Type something"
          value={value}
          onChange={setValue}
          autoComplete="off"
        />
      </Card>
    </Page>
  )
}

function App() {
  return (
    <AppProvider i18n={translations}>
      <AppPage />
    </AppProvider>
  )
}

export default App
